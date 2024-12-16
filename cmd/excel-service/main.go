package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/xuri/excelize/v2"
)

// Struktur untuk menerima pesan
type Message struct {
	CommodityID string `json:"CommodityID"`
	RegionID    string `json:"RegionID"`
}

func createExcelReport(results []domain.PriceHistory, commodityName, regionName string) error {
	// Buat file Excel baru
	f := excelize.NewFile()

	// Buat sheet baru
	sheetName := "Price History Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("error creating sheet: %v", err)
	}
	f.SetActiveSheet(index)

	// Set judul report
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("Price History Report - %s in %s", commodityName, regionName))
	f.MergeCell(sheetName, "A1", "G1")

	// Tulis header
	headers := []string{"No", "Date", "Price", "Unit", "Commodity", "Region"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c3", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style untuk header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#C6EFCE"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating style: %v", err)
	}

	// Apply style ke header
	f.SetCellStyle(sheetName, "A3", "F3", headerStyle)

	// Tulis data
	for i, record := range results {
		row := i + 4 // mulai dari baris ke-4
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), record.CreatedAt.Format("02-01-2006 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), record.Price)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), record.Unit)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), commodityName)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), regionName)
	}

	// Auto-fit column width
	for i := 'A'; i <= 'F'; i++ {
		colName := string(i)
		width, _ := f.GetColWidth(sheetName, colName)
		if width < 15 {
			f.SetColWidth(sheetName, colName, colName, 15)
		}
	}

	// Buat nama file dengan timestamp
	fileName := fmt.Sprintf("./public/reports/price_history_%s_%s_%s.xlsx",
		commodityName,
		regionName,
		time.Now().Format("20060102_150405"))

	// Simpan file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("error saving excel file: %v", err)
	}

	log.Printf("Excel file created successfully: %s", fileName)
	return nil
}

func main() {

	db, err := sql.Open("postgres", "postgres://postgres:123@localhost:5432/go_super_farmer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Koneksi ke RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"report-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg Message
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			// Modifikasi query untuk mengambil data commodity dan region
			currentPrice := domain.Price{
				Commodity: &domain.Commodity{}, // Inisialisasi Commodity
				Region: &domain.Region{ // Inisialisasi Region
					City: &domain.City{}, // Inisialisasi City
				},
			}
			err = db.QueryRow(`
				select p.id, p.commodity_id, p.region_id, p.price, p.unit, p.created_at, co.name, ci.name from prices p join commodities co on p.commodity_id = co.id join regions r on p.region_id = r.id join cities ci on r.city_id = ci.id where p.commodity_id = $1 and p.region_id = $2 and p.deleted_at IS NULL limit 1;
            `, msg.CommodityID, msg.RegionID).Scan(
				&currentPrice.ID,
				&currentPrice.CommodityID,
				&currentPrice.RegionID,
				&currentPrice.Price,
				&currentPrice.Unit,
				&currentPrice.CreatedAt,
				&currentPrice.Commodity.Name,
				&currentPrice.Region.City.Name,
			)
			if err != nil {
				log.Printf("Error querying current price: %v", err)
				continue
			}

			log.Println(currentPrice)

			// Inisialisasi slice results
			var results []domain.PriceHistory

			// Tambahkan current price ke history
			results = append(results, domain.PriceHistory{
				ID:          currentPrice.ID,
				CommodityID: currentPrice.CommodityID,
				RegionID:    currentPrice.RegionID,
				Price:       currentPrice.Price,
				Unit:        currentPrice.Unit,
				CreatedAt:   currentPrice.CreatedAt,
				UpdatedAt:   currentPrice.UpdatedAt,
				DeletedAt:   currentPrice.DeletedAt,
			})

			// Query dan tambahkan price histories
			rows, err := db.Query(`
				select p.id, p.commodity_id, p.region_id, p.price, p.unit, p.created_at, co.name, ci.name from price_histories p join commodities co on p.commodity_id = co.id join regions r on p.region_id = r.id join cities ci on r.city_id = ci.id where p.commodity_id = $1 and p.region_id = $2 and p.deleted_at IS NULL order by p.created_at desc;
            `, msg.CommodityID, msg.RegionID)
			if err != nil {
				log.Printf("Error querying price histories: %v", err)
				continue
			}
			defer rows.Close()

			for rows.Next() {
				history := domain.PriceHistory{
					Commodity: &domain.Commodity{}, // Inisialisasi Commodity
					Region: &domain.Region{ // Inisialisasi Region
						City: &domain.City{}, // Inisialisasi City
					},
				}
				err := rows.Scan(
					&history.ID,
					&history.CommodityID,
					&history.RegionID,
					&history.Price,
					&history.Unit,
					&history.CreatedAt,
					&history.Commodity.Name,
					&history.Region.City.Name,
				)
				if err != nil {
					log.Printf("Error scanning history row: %v", err)
					continue
				}
				results = append(results, history)
			}

			// Buat report Excel
			err = createExcelReport(results, currentPrice.Commodity.Name, currentPrice.Region.City.Name)
			if err != nil {
				log.Printf("Error creating Excel report: %v", err)
				continue
			}

			log.Printf("Successfully processed message and created report")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

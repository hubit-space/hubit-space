package repository

import (
	"fmt"
	"hubit-space/service/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

var monthShort = map[int]string{
	1:  "JAN",
	2:  "FEB",
	3:  "MAR",
	4:  "APR",
	5:  "MAY",
	6:  "JUN",
	7:  "JUL",
	8:  "AUG",
	9:  "SEP",
	10: "OCT",
	11: "NOV",
	12: "DEC",
}

var monthLong = map[int]string{
	1:  "January",
	2:  "February",
	3:  "March",
	4:  "April",
	5:  "May",
	6:  "June",
	7:  "July",
	8:  "August",
	9:  "September",
	10: "October",
	11: "November",
	12: "December",
}

type docNumRepository struct {
	DB *gorm.DB
}

type DocnumRepository interface {
	GetLastDocNumber(docname string, year, month, day int) map[string]any
}

func NewDocnumRepository(db *gorm.DB) *docNumRepository {
	return &docNumRepository{
		DB: db,
	}
}

func (r *docNumRepository) GetLastDocNumber(docname string, year, month, day int) map[string]any {
	defaultDoc, err := r.getDefaultDocFormat(docname)
	if err != nil || defaultDoc == nil {
		return map[string]any{
			"message": "Number of document you requested is not in the settings in the system, please contact your administrator.",
		}
	}

	lastDoc, _ := r.getLastDocFormat(docname, year, month, day, defaultDoc.ResetType)
	var lastnumber int
	var formatStr string

	if lastDoc == nil {
		lastnumber = 1
		if err := r.DB.Table("m_document_number").Create(&model.DocumentNumber{
			DocName:    docname,
			Year:       fmt.Sprintf("%d", year),
			Month:      fmt.Sprintf("%d", month),
			Day:        fmt.Sprintf("%d", day),
			Format:     defaultDoc.Format,
			ResetType:  defaultDoc.ResetType,
			LastNumber: lastnumber,
			IsDefault:  false,
			IsActive:   true,
			CreatedBy:  "SYSTEM",
			CreatedAt:  time.Now(),
			UpdatedBy:  "SYSTEM",
			UpdatedAt:  time.Now(),
		}).Error; err != nil {
			return map[string]any{
				"message": err.Error(),
			}
		}
		formatStr = r.replaceString(defaultDoc.Format, year, month, day, lastnumber)
	} else {
		lastnumber = lastDoc.LastNumber + 1
		if err := r.DB.Table("m_document_number").Where("id = ?", lastDoc.ID).Updates(map[string]any{
			"lastnumber": lastnumber,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return map[string]any{
				"message": err.Error(),
			}
		}
		formatStr = r.replaceString(defaultDoc.Format, year, month, day, lastnumber)
	}

	return map[string]any{
		"format":    defaultDoc.Format,
		"docnum":    formatStr,
		"resettype": defaultDoc.ResetType,
		"counter":   lastnumber,
	}
}

func (r *docNumRepository) getDefaultDocFormat(docname string) (*model.DocumentNumber, error) {
	var doc model.DocumentNumber
	err := r.DB.Table("m_document_number").Where("docname = ? AND is_active = true AND is_default = true", docname).First(&doc).Error
	return &doc, err
}

func (r *docNumRepository) getLastDocFormat(docname string, year, month, day int, resetType string) (*model.DocumentNumber, error) {
	var doc model.DocumentNumber
	query := r.DB.Table("m_document_number").Where("docname = ? AND resettype = ? AND is_default = false", docname, resetType)

	switch resetType {
	case "D":
		query = query.Where("year = ? AND month = ? AND day = ?", year, month, day)
	case "M":
		query = query.Where("year = ? AND month = ?", year, month)
	case "Y":
		query = query.Where("year = ?", year)
	}

	err := query.First(&doc).Error
	if err != nil {
		return nil, err
	}

	return &doc, err
}

func (r *docNumRepository) replaceString(format string, year, month, day, lastnumber int) string {
	number_counter := ""
	counter := strings.Count(format, "X")
	for range counter {
		number_counter = number_counter + "X"
	}

	replace := map[string]string{
		"[DD]":                     fmt.Sprintf("%02d", day),
		"[dd]":                     fmt.Sprintf("%d", day),
		"[MM]":                     fmt.Sprintf("%02d", month),
		"[mm]":                     fmt.Sprintf("%d", month),
		"[MMM]":                    monthShort[month],
		"[MMMM]":                   monthLong[month],
		"[YYYY]":                   fmt.Sprintf("%d", year),
		"[YY]":                     fmt.Sprintf("%02d", year%100),
		"[" + number_counter + "]": fmt.Sprintf("%0*d", counter, lastnumber),
	}

	for k, v := range replace {
		format = strings.ReplaceAll(format, k, v)
	}

	return format
}

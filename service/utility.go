package repository

import (
	"log"

	"hubit-space/service/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OptionRepositoryDB struct {
	DB *gorm.DB
}

type OptionRepository interface {
	GetOptions(c *gin.Context) ([]*model.Option, error)
	GetParamater(parameter string) (*model.Parameter, error)
	GetOptionsInternal(optionName string) ([]*model.Option, error)
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &OptionRepositoryDB{DB: db}
}

func (r *OptionRepositoryDB) GetOptions(c *gin.Context) ([]*model.Option, error) {
	results := []*model.Option{}

	option_name := c.DefaultQuery("option_name", "")

	err := r.DB.
		Table("m_option AS mo").
		Select("mo.id, mod2.id as detail_id, mo.option_name, mo.description, mo.created_at, mod2.option_type, mod2.option_label, mod2.option_value, mod2.is_for_company, mod2.is_mandatory").
		Joins("LEFT JOIN m_option_detail AS mod2 ON mod2.option_type = mo.id").
		Where("mo.option_name = ?", option_name).
		Where("mod2.is_active = true").
		Order("mod2.id ASC").
		Scan(&results).Error

	if err != nil {
		log.Println("Error in get options:", err)
		return nil, err
	}

	return results, nil
}

func (r *OptionRepositoryDB) GetParamater(parameter string) (*model.Parameter, error) {
	var result model.Parameter

	err := r.DB.
		Table("m_parameter AS mp").
		Select("mp.id, mp.parameter_code, mp.parameter_value, mp.description").
		Where("mp.parameter_code = ?", parameter).
		Where("mp.is_active = true").
		First(&result).Error

	if err != nil {
		log.Println("Error in get parameter:", err)
		return nil, err
	}

	return &result, nil
}

func (r *OptionRepositoryDB) GetOptionsInternal(optionName string) ([]*model.Option, error) {
	results := []*model.Option{}

	err := r.DB.
		Table("m_option AS mo").
		Select("mo.id, mod2.id as detail_id, mo.option_name, mo.description, mo.created_at, mod2.option_type, mod2.option_label, mod2.option_value, mod2.is_for_company, mod2.is_mandatory").
		Joins("LEFT JOIN m_option_detail AS mod2 ON mod2.option_type = mo.id").
		Where("mo.option_name = ?", optionName).
		Where("mod2.is_active = true").
		Order("mod2.option_label ASC").
		Scan(&results).Error

	if err != nil {
		log.Println("Error in get options:", err)
		return nil, err
	}

	return results, nil
}

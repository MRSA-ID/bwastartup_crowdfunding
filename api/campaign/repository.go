package campaign

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(Order string, q string) ([]Campaign, error)
	FindByUserID(userID int, Order string, q string) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign) (Campaign, error)
	CreateImage(campaignImage CampaignImage) (CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func (r *repository) FindAll(Order string, q string) ([]Campaign, error){
	var campaigns []Campaign
	var query = "%" + q + "%"
	err := r.db.Order(fmt.Sprintf("id %s", Order)).Where("name LIKE ?",query).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil{
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByUserID(userID int ,Order string, q string) ([]Campaign, error){
	var campaigns []Campaign
	var query = "%" + q + "%"
	err := r.db.Order(fmt.Sprintf("id %s", Order)).Where("user_id = ? AND name LIKE ?", userID, query).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil{
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByID(ID int) (Campaign, error){
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}
	
	return campaign, nil
}

func (r *repository) Save(campaign Campaign) (Campaign, error){
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) Update(campaign Campaign) (Campaign, error){
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) CreateImage(campaignImage CampaignImage) (CampaignImage, error){
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}
	return campaignImage, nil
}

func (r *repository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error){
	// UPDATE campaign_images SET is_primary = false WHERE campaign_id = 1

	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}
	return true, nil
}
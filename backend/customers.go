package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type academicStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type address struct {
	ID              int64     `json:"id"`
	AddressableID   int64     `json:"-"`
	AddressableType string    `json:"-"`
	AddressType     string    `json:"addressType"` // primary or billable_address
	Address1        string    `json:"address1" gorm:"column:address_1"`
	Address2        string    `json:"address2" gorm:"column:address_2"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	PostCode        string    `json:"zip"`
	Country         string    `json:"country"`
	Phone           string    `json:"phone"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type customer struct {
	ID               uint           `json:"id"`
	FirstName        string         `json:"firstName"`
	LastName         string         `json:"lastName"`
	Email            string         `json:"email"`
	AcademicStatusID uint           `json:"academicStatusID"`
	AcademicStatus   academicStatus `gorm:"foreignKey:AcademicStatusID" json:"academicStatus"`
	Addresses        []address      `gorm:"foreignKey:AddressableID" json:"addresses"`
	CreatedAt        time.Time      `json:"createdAt"`
}

func (svc *serviceContext) addOrUpdateCustomer(c *gin.Context) {
	var custReq customer
	err := c.BindJSON(&custReq)
	if err != nil {
		log.Printf("ERROR: invalid customer add/update request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if custReq.ID == 0 {
		log.Printf("INFO: add new customer: %+v", custReq)
		custReq.CreatedAt = time.Now()
		err = svc.DB.Create(&custReq).Error
		if err != nil {
			log.Printf("ERROR: unable to create staff member: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		log.Printf("INFO: update customer: %+v", custReq)
		err = svc.DB.Model(&custReq).Select("LastName", "FirstName", "Email", "AcademicStatusID").Updates(custReq).Error
		if err != nil {
			log.Printf("ERROR: unable to update staff member %d: %s", custReq.ID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	for _, addr := range custReq.Addresses {
		addr.AddressableType = "Customer"
		addr.AddressableID = int64(custReq.ID)
		if addr.ID == 0 {
			addr.CreatedAt = time.Now()
			addr.UpdatedAt = time.Now()
			log.Printf("INFO: add address %+v", addr)
			err = svc.DB.Create(&addr).Error
			if err != nil {
				log.Printf("ERROR: unable to create address: %s", err.Error())
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			log.Printf("INFO: update address %+v", addr)
			err = svc.DB.Model(&addr).Select("Address1", "Address2", "City", "State", "Zip", "Country", "Phone").Updates(addr).Error
			if err != nil {
				log.Printf("ERROR: unable to update staff member %d address: %s", custReq.ID, err.Error())
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	c.JSON(http.StatusOK, custReq)
}

func (svc *serviceContext) getCustomers(c *gin.Context) {
	var customers []customer
	err := svc.DB.Preload("Addresses").Preload("AcademicStatus").Order("last_name asc").Find(&customers).Error
	if err != nil {
		log.Printf("ERROR: unable to get customers: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, customers)
}

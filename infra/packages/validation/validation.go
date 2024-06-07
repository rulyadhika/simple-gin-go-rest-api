package validation

import (
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

func NewValidator() *validator.Validate {
	validation := validator.New()

	// register function to get tag name from json tags.
	validation.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validation.RegisterValidation("ticket_status_custom_validation", func(fl validator.FieldLevel) bool {
		terms := []entity.TicketStatus{entity.TicketStatus_OPEN, entity.TicketStatus_IN_PROGRESS, entity.TicketStatus_RESOLVED, entity.TicketStatus_CLOSED}
		value := fl.Field().Interface().(entity.TicketStatus)

		return slices.Contains(terms, value)
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	err = validation.RegisterValidation("ticket_priority_custom_validation", func(fl validator.FieldLevel) bool {
		terms := []entity.TicketPriority{entity.TicketPriority_LOW, entity.TicketPriority_MED, entity.TicketPriority_HIGH, entity.TicketPriority_CRITICAL}
		value := fl.Field().Interface().(entity.TicketPriority)

		return slices.Contains(terms, value)
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	err = validation.RegisterValidation("user_roles_custom_validation", func(fl validator.FieldLevel) bool {
		terms := []entity.UserType{entity.Role_CLIENT, entity.Role_SUPPORT_AGENT, entity.Role_SUPPORT_SUPERVISOR, entity.Role_ADMINISTRATOR}
		value := fl.Field().Interface().(entity.UserType)

		return slices.Contains(terms, value)
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	return validation
}

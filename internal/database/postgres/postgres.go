package postgres

import (
	"WST_lab6_server/config"
	"WST_lab6_server/internal/database"
	"WST_lab6_server/internal/logging"
	"errors"
	"strconv"
	"strings"

	"WST_lab6_server/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*

 */

type Storage struct {
	DB *gorm.DB
}

/*
Инициализация
*/
func Init() *gorm.DB {
	logging.InitializeLogger()
	var err error
	//Уровень логирования из файла конфигурации
	var logLevel logger.LogLevel
	switch config.GeneralServerSetting.LogLevel {
	case "fatal":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info", "debug":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}
	//Строка подключения к базе данных
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.DatabaseSetting.Host,
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Name,
		config.DatabaseSetting.Port,
		config.DatabaseSetting.SSLMode)
	//Подключаемся к базе данных
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	//Выводим при удачном подключении
	logging.Logger.Info("Database connection established successfully.")
	//Миграция базы данных
	db := conn
	err = db.AutoMigrate(&models.Person{})
	if err != nil {
		log.Fatalf("error creating table: %v", err)
	}
	logging.Logger.Info("Migration completed successfully.")
	//Удаляем таблицу
	db.Exec("DELETE FROM people")
	//Заполняем таблицу из фаила конфигурации
	result := db.Create(&config.GeneralServerSetting.DataSet)
	if result.Error != nil {
		log.Fatalf("error creating table: %v", result.Error)
	}
	//Выводим при удачном заполнениитаблицы
	logging.Logger.Info("Database updated successfully.")
	/*
		//Debug: Запрос к базе и вывод всех данных
	*/
	var results []models.Person
	if err := db.Find(&results).Error; err != nil {
		log.Fatalf("query failed: %v", err)
	}
	for _, record := range results {
		fmt.Println(record)

	}
	fmt.Println("database content in quantity:", len(results), "\n id max:", results[len(results)-1].ID, "id min:", results[0].ID)
	/*
		----
	*/
	//Возвращаем указатель на базу данных
	return db
}

/*
//
Метод поиска в базе данных по запросу
//
*/
func (s *Storage) SearchPerson(searchString string) ([]models.Person, error) {
	var persons []models.Person
	query := s.DB.Model(&models.Person{})
	//Удаляем пробелы из строки поиска
	searchString = strings.TrimSpace(searchString)
	// Проверяем строка является числом, если число ищем по возрасту
	if age, err := strconv.Atoi(searchString); err == nil {
		query = query.Where("age = ?", age)
	} else {
		//Если строка не может быть конвертирована в число ищем по строковым полям
		query = query.Where("name LIKE ? OR surname LIKE ? OR email LIKE ? OR telephone LIKE ?",
			"%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%", "%"+searchString+"%")
	}
	//Выполняем запрос и сохраняем результат в структуру
	if err := query.Find(&persons).Error; err != nil {
		return nil, err
	}
	//Возвращаем результат
	return persons, nil
}

/*
Метод добавления новых данных
*/
func (s *Storage) AddPerson(person *models.Person) (uint, error) {
	//Проверяем наличие записи с таким же email
	if _, err := s.CheckPersonByEmail(person.Email, 0); err == nil {
		return 0, database.ErrEmailExists
	}
	//Создаем запись в базе данных
	if err := s.DB.Create(person).Error; err != nil {
		return 0, err
	}
	return person.ID, nil
}

/*
Метод получения данных
*/
func (s *Storage) GetPerson(id uint) (*models.Person, error) {
	var person models.Person
	//Выполняем запрос к базе данных для получения записи по id
	err := s.DB.First(&person, id).Error
	if err != nil {
		//Возвращаем ошибку при выполнении запроса к базе данных
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrPersonNotFound
		}
		return nil, err
	}
	//Возвращаем результат
	return &person, nil
}

/*
Метод обновления данных по id
*/
func (s *Storage) UpdatePerson(person *models.Person) error {
	//Выполняем запрос к базе данных для обновления записи
	result := s.DB.Model(&models.Person{}).Where("id = ?", person.ID).Updates(models.Person{
		Name:      person.Name,
		Surname:   person.Surname,
		Age:       person.Age,
		Email:     person.Email,
		Telephone: person.Telephone,
	})

	if result.Error != nil {
		//Возвращаем ошибку при выполнении запроса к базе данных
		return result.Error
	}

	if result.RowsAffected == 0 {
		//Возвращаем ошибку если запись не найдена для обновления
		return database.ErrPersonNotFound
	}
	//Возвращаем ничего при успехе
	return nil
}

/*
Метод удаления данных по id
*/
func (s *Storage) DeletePerson(person *models.Person) error {
	//Выполняем запрос к базе данных для удаления записи по id
	result := s.DB.Delete(&person)
	if result.Error != nil {
		//Возвращаем ошибку при выполнении запроса к базе данных
		return result.Error
	}
	//Возвращаем ошибку при выполнении запроса к базе данных
	return result.Error
}

/*
Метод получения всех данных
*/
func (s *Storage) GetAllPersons() ([]models.Person, error) {
	var persons []models.Person
	//Выполняем запрос к базе данных для получения всех записей
	err := s.DB.Find(&persons).Error
	if err != nil {
		return nil, err
	}
	//Возвращаем результат
	return persons, nil
}

/*
Метод проверки наличия записи по email
*/
func (s *Storage) CheckPersonByEmail(email string, excludeId uint) (*models.Person, error) {
	var person models.Person
	// Выполняем запрос к базе данных для поиска по email
	if err := s.DB.Where("email = ? AND id != ?", email, excludeId).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//Возвращаем кастомную ошибку (Запись не найдена)
			return nil, database.ErrPersonNotFound
		}
		//Возвращаем ошибку
		return nil, err
	}
	//Возвращаем запись
	return &person, nil
}

/*
Метод проверки наличия записи по id
*/
func (s *Storage) CheckPersonByIDHandler(id uint) (bool, error) {
	var person models.Person
	//Выполняем запрос к базе данных для поиска по id
	result := s.DB.First(&person, id)
	if result.Error != nil {
		//Проверяем наличие записи по id
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("Record not found")
		} else {
			fmt.Println("Error when executing the request:", result.Error)
		}
	} else {
		fmt.Println("The record was found with CheckPersonByIDHandler:", person)
		return true, nil
	}
	//Возвращаем false при успехе
	return false, nil
}

package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000 // количество метров в одном километре
	MinInHours = 60   // количество минут в одном часе
	LenStep    = 0.65 // длина одного шага
	CmInM      = 100  // количество сантиметров в одном метре
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // тип тренировки
	Action       int           // количество повторов(шаги, гребки при плавании)
	LenStep      float64       // длина одного шага или гребка в м
	Duration     time.Duration // продолжительность тренировки
	Weight       float64       // вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	if t.Duration.Hours() == 0 {
		return 0
	}
	return t.distance() / t.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий на тренировке.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string        // тип тренировки
	Duration     time.Duration // длительность тренировки
	Distance     float64       // расстояние, которое преодолел пользователь
	Speed        float64       // средняя скорость, с которой двигался пользователь
	Calories     float64       // количество потраченных килокалорий на тренировке
}

// TrainingInfo возвращает структуру InfoMessage с информацией о тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для всех типов тренировок.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для бега
const (
	CaloriesMeanSpeedMultiplier = 18
	CaloriesMeanSpeedShift      = 1.79
)

// Running структура для бега
type Running struct {
	Training
}

func (r Running) Calories() float64 {
	if r.Duration.Hours() == 0 {
		return 0
	}
	speed := r.meanSpeed()
	return (CaloriesMeanSpeedMultiplier*speed + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours
}

func (r Running) TrainingInfo() InfoMessage {
	info := r.Training.TrainingInfo()
	info.Calories = r.Calories()
	return info
}

// Константы для ходьбы
const (
	CaloriesWeightMultiplier      = 0.035
	CaloriesSpeedHeightMultiplier = 0.029
	KmHInMsec                     = 0.278
)

// Walking структура для ходьбы
type Walking struct {
	Training
	Height float64
}

func (w Walking) Calories() float64 {
	if w.Duration.Hours() == 0 || w.Height == 0 {
		return 0
	}
	speedMsec := w.meanSpeed() * KmHInMsec
	heightM := w.Height / CmInM
	return (CaloriesWeightMultiplier*w.Weight + math.Pow(speedMsec, 2)/heightM*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours
}

func (w Walking) TrainingInfo() InfoMessage {
	info := w.Training.TrainingInfo()
	info.Calories = w.Calories()
	return info
}

// Константы для плавания
const (
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

// Swimming структура для плавания
type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

func (s Swimming) meanSpeed() float64 {
	if s.Duration.Hours() == 0 {
		return 0
	}
	return float64(s.LengthPool*s.CountPool) / MInKm / s.Duration.Hours()
}

func (s Swimming) Calories() float64 {
	if s.Duration.Hours() == 0 {
		return 0
	}
	speed := s.meanSpeed()
	return (speed + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     s.distance(),
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// ReadData обрабатывает данные тренировки
func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	return fmt.Sprint(info)
}

func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))
}

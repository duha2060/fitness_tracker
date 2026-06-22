package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

// validateInput  возвращает ошибки
func validateInput(steps int, weight, height float64, duration time.Duration) error {
	if weight <= 0 || height <= 0 || steps <= 0 || duration <= 0 {
		return errors.New("некорректные входные данные")
	}
	return nil
}

func parseTraining(data string) (int, string, time.Duration, error) {
	parseData := strings.Split(data, ",")

	if len(parseData) != 3 {
		return 0, "", 0, errors.New("invalid data format")
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(parseData[2])
	if err != nil {
		return 0, "", 0, err
	}

	// Тесты ожидают ошибку при парсинге, если значения <= 0
	if steps <= 0 || duration <= 0 {
		return 0, "", 0, errors.New("invalid steps or duration")
	}

	return steps, parseData[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return stepLengthCoefficient * height * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / (duration.Minutes() / minInH)
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	switch activityType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		err = errors.New("неизвестный тип тренировки")
		log.Println(err)
		return "", err
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / minInH * walkingCaloriesCoefficient, nil
}

package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Неправильный формат данных
var errorFormateData = errors.New("wrong data format")
var errorUnknownActivity = errors.New("неизвестный тип тренировки")

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Парсит строку и возвращает её составные части с нужными типами
// Принимает на вход строку (data) в формате: <кол-во шагов>,<вид активности>,<длительность активности>
// Возвращает:
//   - количество шагов (int)
//   - вид активности (string)
//   - продолжительность активности (time.Duration)
//   - ошибка, если что-то пошло не так. (error)
func parseTraining(data string) (int, string, time.Duration, error) {
	parseData := strings.Split(data, ",")
	if len(parseData) != 3 {
		return 0, "", 0, fmt.Errorf("%w: three comma-separated values are expected", errorFormateData)
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w. %w", errorFormateData, err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("%w: first value in data must be > 0", errorFormateData)
	}

	if parseData[1] == "" {
		return 0, "", 0, fmt.Errorf("%w: second value in data is not defined", errorFormateData)
	}

	duration, err := time.ParseDuration(parseData[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w: %w", errorFormateData, err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("%w: third value in the data must be > 0", errorFormateData)
	}

	return steps, parseData[1], duration, nil
}

// Возвращает дистанцию в километрах
// Принимает на вход:
//
//	steps  - количество шагов
//	height - рост пользователя
func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	return (stepLen * float64(steps)) / mInKm
}

// Возвращает среднюю скорость в км/ч
// Принимает на вход:
//
//	steps    - количество шагов
//	height   - рост пользователя
//	duration - продолжительность активности
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	return distance / duration.Hours()
}

// Возвращает информацию о тренировке или ошибку
// Принимает на вход:
//
//	data   - <кол-во шагов>,<вид активности>,<длительность активности>
//	weight - вес пользователя в килограммах
//	height - рост пользователя в метрах
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)
	var calories float64 = 0

	switch strings.ToUpper(activity) {
	case "ХОДЬБА":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "БЕГ":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		err = fmt.Errorf("%w. %w", errorFormateData, errorUnknownActivity)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), distance, meanSpeed, calories), nil
}

// Возвращает количество калорий, потраченных при беге или ошибку.
// Принимает на вход:
//
//	steps    - количество шагов
//	weight   - вес пользователя
//	height   - рост пользователя
//	duration - продолжительность активности
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("%w: steps must be > 0", errorFormateData)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("%w: weight must be > 0", errorFormateData)
	}
	if height <= 0 {
		return 0, fmt.Errorf("%w: height must be > 0", errorFormateData)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%w: duration must be > 0", errorFormateData)
	}

	meanSpeed := meanSpeed(steps, height, duration)
	if meanSpeed <= 0 {
		return 0, fmt.Errorf("%w: error in calculating average speed", errorFormateData)
	}

	durationInMinutes := duration.Minutes()
	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

// Возвращает количество калорий, потраченных при ходьбе или ошибку.
// Принимает на вход:
//
//	steps    - количество шагов
//	weight   - вес пользователя
//	height   - рост пользователя
//	duration - продолжительность активности
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return runningSpentCalories * walkingCaloriesCoefficient, nil
}

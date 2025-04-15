package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

// Неправильный формат данных
var errorFormateData = errors.New("wrong data format")

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// Парсит строку и возвращает её составные части с нужными типами
// Принимает на вход строку (data) в формате: <кол-во шагов>,<длительность прогулки>
// Возвращает:
//   - количество шагов (int)
//   - длительность прогулки (time.Duration)
//   - ошибка, если что-то пошло не так. (error)
func parsePackage(data string) (int, time.Duration, error) {
	parseData := strings.Split(data, ",")
	if len(parseData) != 2 {
		return 0, 0, fmt.Errorf("%w: two comma-separated values are expected", errorFormateData)
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", errorFormateData, err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("%w: first value in data must be > 0", errorFormateData)
	}

	duration, err := time.ParseDuration(parseData[1])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", errorFormateData, err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("%w: second value in data must be > 0", errorFormateData)
	}

	return steps, duration, nil
}

// Возвращает информацию об активности за день
// Принимает на вход:
//
//	data — строка с данными в формате: <кол-во шагов>,<продолжительность прогулки в формате 3h50m (3 часа 50 минут)>
//	weight — вес пользователя в килограммах.
//	height — рост пользователя в метрах.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distanceKm := (float64(steps) * stepLength) / mInKm

	caloriesDay, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, caloriesDay)
}

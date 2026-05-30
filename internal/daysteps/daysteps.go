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

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parseData := strings.Split(data, ",")

	if len(parseData) != 2 {
		return 0, 0, errors.New("invalid data format")
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, 0, err
	}

	duration, err := time.ParseDuration(parseData[1])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 || duration <= 0 {
		return 0, 0, errors.New("invalid steps or duration")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distance := float64(steps) * stepLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distance,
		calories,
	)
}

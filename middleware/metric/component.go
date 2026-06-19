package metric

import (
	"github.com/Compogo/compogo"
)

var (
	// RequestCountComponent — компонент для сбора метрики количества запросов.
	// Считает количество HTTP-запросов с разбивкой по эндпоинтам и кодам ответа.
	RequestCountComponent = &compogo.Component{
		Init: compogo.StepFunc(func(container compogo.Container) error {
			return container.Provide(NewRequestCount)
		}),
	}

	// DurationComponent — компонент для сбора метрики длительности запросов.
	// Измеряет время выполнения HTTP-запросов с разбивкой по эндпоинтам.
	DurationComponent = &compogo.Component{
		Init: compogo.StepFunc(func(container compogo.Container) error {
			return container.Provide(NewDuration)
		}),
	}
)

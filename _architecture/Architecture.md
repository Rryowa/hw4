### Архитектурный паттерн MVC
• Модель предоставляет данные и реагирует на команды контроллера, изменяя своё состояние.  
• Представление отвечает за отображение данных модели пользователю, реагируя на изменение модели.  
• Контроллер(Сервис) интерпретирует действия пользователя, оповещая модель о необходимости изменений.

### Strategy behavioral pattern
https://refactoring.guru/design-patterns/strategy/go/example
```go
type packageContext struct {
    strategies map[models.PackageType]PackageStrategy
}

type PackageService interface {
    ValidatePackage(weight models.Weight, packageType models.PackageType) error
    ApplyPackage(order *models.Order, packageType models.PackageType)
}

type PackageStrategy interface {
    Validate(weight models.Weight) error
    Apply(order *models.Order)
}

func NewPackageService() PackageService {
    return &packageContext{
        strategies: map[models.PackageType]PackageStrategy{
            FilmType:   NewFilmPackage(),
            PacketType: NewPacketPackage(),
            BoxType:    NewBoxPackage(),
        },
    }
}

func (pc *packageContext) ValidatePackage(weight models.Weight, packageType models.PackageType) error {
    if strategy, ok := pc.strategies[packageType]; ok {
        return strategy.Validate(weight)
    }
	
    return util.ErrPackageTypeInvalid
}

func (pc *packageContext) ApplyPackage(order *models.Order, packageType models.PackageType) {
    if strategy, ok := pc.strategies[packageType]; ok {
        strategy.Apply(order)
        return
    }
    //Assume that film has no weight limit
    pc.strategies[FilmType].Apply(order)
}

func NewBoxPackage() *BoxPackage {
    return &BoxPackage{}
}

func (p *BoxPackage) Validate(weight models.Weight) error {
    if weight < MaxBoxWeight {
        return nil
	}
    return util.ErrWeightExceeds
}

func (c *BoxPackage) Apply(order *models.Order) {
    order.PackageType = BoxType
    order.PackagePrice = BoxPrice
    order.OrderPrice += BoxPrice
}
```

## Используемые стандарты описания архитектуры
Для описания архитектуры я использовал *диаграмму классов* и  
*диаграмму последовательности* стандарта UML,

### Почему UML
UML я выбрал, потому что C4 Model непригодна для настолько маленького проекта,  
2 и 3 слои выглядят идентично из-за отсутствия внешнего сервиса.


### Почему диаграмма последовательности
Я выбрал диаграмму последовательности, так как с её помощью можно понять, как пользователь  
и программа взаимодействуют от начал до конца.  
Некоторого рода общая картина взаимодействия.


### Почему диаграмма классов
Я выбрал диаграмму классов, потому что таким образом можно углубиться в то как работает программа,  
при этом не углубляясь в технические детали реализации методов интерфейса.
Нечто среднее между общей картиной и точной реализацией каждого интерфейса, если коллега захочет  
углубиться в детали реализации, то они будут доступны в самом коде проекта.

### Почему в диаграмме классов есть пакеты(namespace)?
Потому что пакеты в гошке это круто и я думаю, это не усложняет процесс понимания, а наоборот помогает.  
Но я могу и ошибаться...

### Таким образом, коллега сможет при анализе диаграммы классов держать в голове общую картину программы и при этом разбираться в нужном пакете или классе(структуре) если у него возникнут вопросы!
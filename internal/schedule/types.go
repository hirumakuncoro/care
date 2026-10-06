package schedule

type SlotKind string

const (
	KindConfirm   SlotKind = "confirm"
	KindKickCount SlotKind = "kickcount"
)

type ItemCategory string

const (
	CatMeal       ItemCategory = "meal"
	CatSnack      ItemCategory = "snack"
	CatFruit      ItemCategory = "fruit"
	CatSupplement ItemCategory = "supplement"
	CatMilk       ItemCategory = "milk"
	CatActivity   ItemCategory = "activity"
)

type ClockTime struct {
	Hour int
	Min  int
}

type Item struct {
	Name     string
	Category ItemCategory
	Note     string
}

type Slot struct {
	ID       string
	At       ClockTime
	Title    string
	Items    []Item
	Kind     SlotKind
	WindowMs int64
}

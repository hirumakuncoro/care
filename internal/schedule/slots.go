package schedule

var Slots = []*Slot{
	{
		ID:       "breakfast",
		At:       ClockTime{Hour: 7, Min: 0},
		Title:    "Sarapan",
		Kind:     KindConfirm,
		WindowMs: 60 * 60 * 1000, // 60 menit
		Items: []Item{
			{Name: "Maltofer #1", Category: CatSupplement, Note: "sesudah makan"},
			{Name: "Buah vit C", Category: CatFruit},
		},
	},
	{
		ID:       "fruit",
		At:       ClockTime{Hour: 10, Min: 0},
		Title:    "Buah & Susu",
		Kind:     KindConfirm,
		WindowMs: 60 * 60 * 1000,
		Items: []Item{
			{Name: "Buah (alpukat/pisang)", Category: CatFruit},
			{Name: "Susu Diamond #1", Category: CatMilk},
		},
	},
	{
		ID:       "lunch",
		At:       ClockTime{Hour: 13, Min: 0},
		Title:    "Makan Siang",
		Kind:     KindConfirm,
		WindowMs: 60 * 60 * 1000,
		Items: []Item{
			{Name: "Hi-Bone + D3 + I-Folic", Category: CatSupplement, Note: "sesudah makan"},
		},
},
	{
		ID:       "snack_pro",
		At:       ClockTime{Hour: 16, Min: 0},
		Title:    "Camilan Protein",
		Kind:     KindConfirm,
		WindowMs: 60 * 60 * 1000,
		Items: []Item{
			{Name: "Telur rebus/kacang", Category: CatSnack},
			{Name: "Susu Diamond #2", Category: CatMilk},
		},
	},
	{
		ID:       "dinner",
		At:       ClockTime{Hour: 19, Min: 0},
		Title:    "Makan Malam",
		Kind:     KindConfirm,
		WindowMs: 60 * 60 * 1000,
		Items: []Item{
			{Name: "Maltofer #2", Category: CatSupplement, Note: "sesudah makan"},
			{Name: "Buah vit C", Category: CatFruit},
		},
	},
	{
		ID:       "kickcount",
		At:       ClockTime{Hour: 20, Min: 0},
		Title:    "Ritual Kick Count",
		Kind:     KindKickCount,
		WindowMs: 90 * 60 * 1000, // 90 menit
		Items: []Item{
			{Name: "10 kick count", Category: CatActivity},
		},
},
	{
		ID:       "snack_cal",
		At:       ClockTime{Hour: 21, Min: 30},
		Title:    "Camilan Kalori",
		Kind:     KindConfirm,
		WindowMs: 60 * 60 * 1000,
		Items: []Item{
			{Name: "Biskuit gandum + pisang", Category: CatSnack},
			{Name: "Susu Diamond #3", Category: CatMilk},
		},
	},
}

func SlotByID(id string) *Slot {
	for _, slot := range Slots {
		if slot.ID == id {
			return slot
		}
	}
	return nil
}
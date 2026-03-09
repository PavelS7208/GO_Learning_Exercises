package person

import "GamePerson/internal/model/game/creatures/base/entity"

// Компилятор проверит соответствие интерфейсам
var (
	_ entity.Combatant     = (*person)(nil)
	_ entity.Wealthy       = (*person)(nil)
	_ entity.Magical       = (*person)(nil)
	_ entity.Experienced   = (*person)(nil)
	_ entity.Reputable     = (*person)(nil)
	_ entity.PropertyOwner = (*person)(nil)
	_ entity.FamilyMember  = (*person)(nil)
	_ entity.Validatable   = (*person)(nil) // ← критически важно
)

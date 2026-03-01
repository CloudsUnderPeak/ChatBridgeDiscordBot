package gamble

import (
	pkgSql "discord-chatbot/pkg/sql"
)

func CreateDbGamer(gamer *Gamer) error {
	gamer.dbMu.Lock()
	defer gamer.dbMu.Unlock()
	gamerData := &pkgSql.DiscordGambleGamer{
		Id:    gamer.id,
		Name:  gamer.name,
		Chips: gamer.chips,
	}
	return pkgSql.Database.Create(gamerData).Error
}

func ReadDbGamer(id string) (*Gamer, error) {
	var dbGamer = pkgSql.DiscordGambleGamer{}
	if err := pkgSql.Database.First(&dbGamer, "Id = ?", id).Error; err != nil {
		return nil, err
	}
	var gamer = &Gamer{
		id:    dbGamer.Id,
		name:  dbGamer.Name,
		chips: dbGamer.Chips,
	}
	return gamer, nil
}

func ReadAllDbGamer() ([]*Gamer, error) {
	var dbGamers = []pkgSql.DiscordGambleGamer{}
	if err := pkgSql.Database.Find(&dbGamers).Error; err != nil {
		return nil, err
	}
	var gamers []*Gamer
	for _, dbGamer := range dbGamers {
		var gamer = &Gamer{
			id:    dbGamer.Id,
			name:  dbGamer.Name,
			chips: dbGamer.Chips,
		}
		gamers = append(gamers, gamer)
	}
	return gamers, nil
}

func (g *Gamer) UpdateDbChips() error {
	g.dbMu.Lock()
	defer g.dbMu.Unlock()
	return pkgSql.Database.Model(&pkgSql.DiscordGambleGamer{}).
		Where("id = ?", g.id).
		Update("Chips", g.chips).Error
}

func (g *Gamer) UpdateDbName() error {
	g.dbMu.Lock()
	defer g.dbMu.Unlock()
	return pkgSql.Database.Model(&pkgSql.DiscordGambleGamer{}).
		Where("id = ?", g.id).
		Update("Name", g.name).Error
}

func DeleteDbGamer(gamer *Gamer) error {
	gamer.dbMu.Lock()
	defer gamer.dbMu.Unlock()
	return pkgSql.Database.Delete(&pkgSql.DiscordGambleGamer{}, "Id = ?", gamer.id).Error
}

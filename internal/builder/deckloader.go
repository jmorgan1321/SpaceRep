package builder

import "github.com/jmorgan1321/SpaceRep/internal/core"

type deckLoader interface {
	LoadDeck([]string) (*core.Deck, error)
	SaveDeck()
}

type fromDiskDeckLoader struct {
}

func (m *fromDiskDeckLoader) LoadDeck(sets []string) (*core.Deck, error) {
	deck := &core.Deck{}
	// for _, set := range sets {
	//     root := b.rootPath() + set
	//     info, err := getDeckInfo(root + "/cards/cards.info")
	//     if err != nil {
	//         return nil, err
	//     }

	//     tmpls, err := getCardTemplatesFromDisk(root+"/cards", info.Display)
	//     if err != nil {
	//         return nil, err
	//     }

	//     cards := makeCards(set, info.Info, tmpls)
	//     updateBuckets(cards)

	//     SaveDeck(root+"/cards/cards.info", info.Display, cards)

	//     for _, c := range cards {
	//         deck[c.Bucket()] = append(deck[c.Bucket()], c)
	//     }
	// }

	return deck, nil
}
func (m *fromDiskDeckLoader) SaveDeck() {
	// info := []*core.Info{}
	// for _, c := range cards {
	// 	// TODO: remove hard coded
	// 	info = append(info, c.Stats())
	// }

	// d := deckInfo{Display: display, Info: info}
	// b, _ := json.MarshalIndent(d, "", "\t")

	// f, err := os.Create(path)
	// if err != nil {
	// 	panic("file read error with: " + path)
	// }
	// defer f.Close()

	// if _, err := f.Write(b); err != nil {
	// 	panic("saving " + path + ":" + err.Error())
	// }
}

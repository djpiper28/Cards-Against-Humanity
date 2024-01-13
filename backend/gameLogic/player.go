package gameLogic

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type Player struct {
	Id          uuid.UUID
	Name        string
	Hand        map[int]*WhiteCard
	CurrentPlay []*WhiteCard
	Connected   bool
	Points      int
}

const (
	MaxPlayerNameLength = 20
	MinPlayerNameLength = 3
)

func NewPlayer(Name string) (*Player, error) {
	if len(Name) > MaxPlayerNameLength || len(Name) < MinPlayerNameLength {
		return nil, errors.New(fmt.Sprintf("Length of name must be between %d and %d (exclusive exclusive)", MinPlayerNameLength, MaxPlayerNameLength))
	}

	return &Player{Id: uuid.New(),
		Name:      Name,
		Hand:      make(map[int]*WhiteCard),
		Connected: true}, nil
}

func (p *Player) hasCard(card *WhiteCard) bool {
	_, found := p.Hand[card.Id]
	return found
}

func (p *Player) PlayCard(cards []*WhiteCard) error {
	if cards == nil {
		return errors.New("Cannot play nil cards")
	}

	if p.CurrentPlay != nil {
		return errors.New("Cards have already been played")
	}

	cardsSeen := make(map[int]bool)
	for _, card := range cards {
		_, found := cardsSeen[card.Id]
		if found {
			return errors.New("Card is in your play more than once")
		}
		cardsSeen[card.Id] = true

		if !p.hasCard(card) {
			return errors.New("Cannot find the card in the hand")
		}
	}

	cardsCopy := make([]*WhiteCard, len(cards))
	copy(cardsCopy, cards)

	p.CurrentPlay = cardsCopy

	for _, card := range cards {
		delete(p.Hand, card.Id)
	}
	return nil
}

func (p *Player) AddCardToHand(card *WhiteCard) error {
	if p.hasCard(card) {
		msg := "Cannot add duplicate cards to the hand"
		log.Println(msg)
		return errors.New(msg)
	}

	p.Hand[card.Id] = card
	return nil
}

// Used in tests
func (p *Player) CardsInHand() int {
	count := 0
	for range p.Hand {
		count++
	}
	return count
}

func (p *Player) FinaliseRound() {
	p.CurrentPlay = nil
}

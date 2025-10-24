package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/list"
)

const (
	pattyLimit   int = 2
	sauceLimit   int = 2
	garnishLimit int = 4
)

var (
	errTooManyPatties = errors.New("too many patties")
	errTooMuchSauce   = errors.New("too much sauce")
	errTooMuchGarnish = errors.New("too much garnish")
)

type Burger struct {
	Name string
	Bun
	Patties   map[Patty]int
	Sauces    map[Sauce]int
	Garnishes map[Garnish]int
}

func (b Burger) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Burger: %s\n", b.Name))
	sb.WriteString(fmt.Sprintf("Bun: %s\n", b.Bun))

	items := make([]string, 0, len(b.Patties))
	for patty, count := range b.Patties {
		items = append(items, fmt.Sprintf("%d %s", count, patty))
	}
	sb.WriteString(fmt.Sprintf("Patty: %s\n", strings.Join(items, ", ")))

	items = make([]string, 0, len(b.Sauces))
	for sauce, count := range b.Sauces {
		items = append(items, fmt.Sprintf("%d %s", count, sauce))
	}
	sb.WriteString(fmt.Sprintf("Sauce: %s\n", strings.Join(items, ", ")))

	items = make([]string, 0, len(b.Garnishes))
	for garnish, count := range b.Garnishes {
		items = append(items, fmt.Sprintf("%d %s", count, garnish))
	}
	sb.WriteString(fmt.Sprintf("Garnish: %s\n", strings.Join(items, ", ")))

	return sb.String()
}

func (b Burger) PattyCount() int {
	return list.Sum(dict.Values(b.Patties))
}

func (b Burger) SauceCount() int {
	return list.Sum(dict.Values(b.Sauces))
}

func (b Burger) GarnishCount() int {
	return list.Sum(dict.Values(b.Garnishes))
}

type OptionFn[T any] func(*T) error

type (
	Bun     string
	Patty   string
	Sauce   string
	Garnish string
)

var (
	SmallBun  Bun = "small"
	MediumBun Bun = "medium"
	LargeBun  Bun = "large"
)

var (
	BeefPatty    Patty = "beef"
	ChickenPatty Patty = "chicken"
	VeganPatty   Patty = "vegan"
)

var (
	Ketchup Sauce = "ketchup"
	Mustard Sauce = "mustard"
	Mayo    Sauce = "mayo"
)

var (
	Cheese  Garnish = "cheese"
	Tomato  Garnish = "tomato"
	Lettuce Garnish = "lettuce"
	Pickle  Garnish = "pickle"
	Onion   Garnish = "onion"
)

func NewBurger(options ...OptionFn[Burger]) (*Burger, error) {
	// Default burger
	burger := &Burger{
		Name:      "burger",
		Bun:       MediumBun,
		Patties:   make(map[Patty]int),
		Sauces:    make(map[Sauce]int),
		Garnishes: make(map[Garnish]int),
	}

	// Decorate with options
	for _, opt := range options {
		if err := opt(burger); err != nil {
			return nil, err
		}
	}

	return burger, nil
}

func WithName(name string) OptionFn[Burger] {
	return func(burger *Burger) error {
		burger.Name = name
		return nil
	}
}

func WithBun(bun Bun) OptionFn[Burger] {
	return func(burger *Burger) error {
		burger.Bun = bun
		return nil
	}
}

func WithPatty(patty Patty) OptionFn[Burger] {
	return WithMultiPatty(patty, 1)
}

func WithMultiPatty(patty Patty, qty int) OptionFn[Burger] {
	return func(burger *Burger) error {
		if burger.PattyCount()+qty > pattyLimit {
			return errTooManyPatties
		}
		burger.Patties[patty] += qty
		return nil
	}
}

func WithSauce(sauce Sauce) OptionFn[Burger] {
	return WithMultiSauce(sauce, 1)
}

func WithMultiSauce(sauce Sauce, qty int) OptionFn[Burger] {
	return func(burger *Burger) error {
		if burger.SauceCount()+qty > sauceLimit {
			return errTooMuchSauce
		}
		burger.Sauces[sauce] += qty
		return nil
	}
}

func WithGarnish(garnish Garnish) OptionFn[Burger] {
	return WithMultiGarnish(garnish, 1)
}

func WithMultiGarnish(garnish Garnish, qty int) OptionFn[Burger] {
	return func(burger *Burger) error {
		if burger.GarnishCount()+qty > garnishLimit {
			return errTooMuchGarnish
		}
		burger.Garnishes[garnish] += qty
		return nil
	}
}

func Test2() {
	burger0, err := NewBurger()
	displayBurger(burger0, err)

	burger1, err := NewBurger(
		WithName("big mac"),
		WithBun(LargeBun),
		WithMultiPatty(BeefPatty, 2),
		WithSauce(Ketchup),
		WithGarnish(Cheese),
		WithGarnish(Tomato),
		WithGarnish(Onion),
		WithGarnish(Pickle),
	)
	displayBurger(burger1, err)

	burger2, err := NewBurger(
		WithName("monster burger"),
		WithBun(LargeBun),
		WithPatty(BeefPatty),
		WithMultiPatty(ChickenPatty, 2),
	)
	displayBurger(burger2, err)

	burger3, err := NewBurger(
		WithName("saucy burger"),
		WithSauce(Ketchup),
		WithSauce(Mayo),
		WithSauce(Mustard),
	)
	displayBurger(burger3, err)

	burger4, err := NewBurger(
		WithName("garnished burger"),
		WithMultiGarnish(Cheese, 2),
		WithMultiGarnish(Tomato, 2),
		WithGarnish(Lettuce),
	)
	displayBurger(burger4, err)
}

func displayBurger(burger *Burger, err error) {
	if err == nil {
		fmt.Println(burger)
	} else {
		fmt.Printf("Error: %s\n\n", err.Error())
	}
}

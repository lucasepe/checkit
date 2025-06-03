package parser_test

import (
	"strings"
	"testing"

	"github.com/lucasepe/checkit/internal/parser"
)

func TestParser(t *testing.T) {
	input := `
# Lasciare un appartamento venduto

## Utenze

- Disdetta o voltura utenze
> Quando farla: almeno 10-15 giorni prima del rogito.
> Documenti: codice cliente, letture contatori, documento d'identità.

- Lettura contatori
> Falla il giorno della consegna chiavi.

## TARI

* Comunicazione cessazione
> Entro 30 giorni dal rogito.

* Nuova iscrizione
`
	lst, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Errore durante il parsing: %v", err)
	}

	if lst.Title != "Lasciare un appartamento venduto" {
		t.Errorf("Titolo checklist errato. Atteso: %q, Ottenuto: %q", "Lasciare un appartamento venduto", lst.Title)
	}

	if len(lst.Groups) != 2 {
		t.Errorf("Numero gruppi errato. Atteso: 2, Ottenuto: %d", len(lst.Groups))
	}

	utenze := lst.Groups[0]
	if utenze.Title != "Utenze" {
		t.Errorf("Titolo primo gruppo errato. Atteso: %q, Ottenuto: %q", "Utenze", utenze.Title)
	}
	if len(utenze.Items) != 2 {
		t.Errorf("Numero item in 'Utenze' errato. Atteso: 2, Ottenuto: %d", len(utenze.Items))
	}
	if utenze.Items[0].Title != "Disdetta o voltura utenze" {
		t.Errorf("Titolo primo item errato. Atteso: %q, Ottenuto: %q", "Disdetta o voltura utenze", utenze.Items[0].Title)
	}
	if len(utenze.Items[0].Notes) != 2 {
		t.Errorf("Note primo item errate. Atteso: 2, Ottenuto: %d", len(utenze.Items[0].Notes))
	}

	tari := lst.Groups[1]
	if tari.Title != "TARI" {
		t.Errorf("Titolo secondo gruppo errato. Atteso: %q, Ottenuto: %q", "TARI", tari.Title)
	}
	if len(tari.Items) != 2 {
		t.Errorf("Numero item in 'TARI' errato. Atteso: 2, Ottenuto: %d", len(tari.Items))
	}
	if tari.Items[1].Title != "Nuova iscrizione" {
		t.Errorf("Titolo secondo item in 'TARI' errato. Atteso: %q, Ottenuto: %q", "Nuova iscrizione", tari.Items[1].Title)
	}
}

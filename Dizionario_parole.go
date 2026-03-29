package main

import (
	"bufio"
	. "fmt"
	"os"
	"strings"
	"unicode"
)

type dizionario struct {
	sequenze *map[byte][]tipoSequenza
	indice   *map[int][]string
}

type tipoSequenza struct {
	sequenza string
	isWord   bool
}

func main() {
	d := newDizionario()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		d = esegui(d, scanner.Text())
	}
}

func newDizionario() dizionario {
	d := dizionario{}
	d.crea()
	return d
}

func (d *dizionario) crea() {
	*d = dizionario{
		sequenze: &map[byte][]tipoSequenza{},
		indice:   &map[int][]string{},
	}
}

func esegui(d dizionario, s string) dizionario {
	dati := strings.Fields(s)
	if len(dati) == 0 {
		return d
	}

	switch dati[0] {
	case "c":
		switch len(dati) {
		case 1:
			newDizionario()
			*d.sequenze = make(map[byte][]tipoSequenza)
		case 2:
			newDizionario()
			*d.sequenze = make(map[byte][]tipoSequenza)
			*d.indice = make(map[int][]string)
			d.carica(dati[1])
		case 3:
			d.catena(dati[1], dati[2])
		}
	case "p":
		d.stampa_parole()
	case "s":
		d.stampa_schemi()
	case "i":
		if len(dati) == 2 {
			d.inserisci(dati[1])
		}
	case "e":
		if len(dati) == 2 {
			d.elimina(dati[1])
		}
	case "r":
		if len(dati) == 2 {
			d.ricerca(dati[1])
		}
	case "d":
		if len(dati) == 3 {
			Println(distanza(dati[1], dati[2]))
		}
	case "t":
		os.Exit(0)
	}
	return d
}

func creaChiaveMinuscola(c rune) byte {
	return byte(unicode.ToLower(c))
}

func checkIniziale(w string) (bool, byte) {
	iniziale := byte(unicode.ToLower(rune(w[0])))
	isSchema := false
	for _, r := range w {
		if unicode.IsUpper(r) {
			isSchema = true
			iniziale = creaChiaveMinuscola(r)
			break
		}
	}
	return isSchema, iniziale
}

func (d *dizionario) checkDuplicati(w string, iniziale byte) bool {
	for _, presente := range (*d.sequenze)[iniziale] {
		if presente.sequenza == w {
			return true
		}
	}
	return false
}

func (d *dizionario) carica(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, w := range strings.Fields(scanner.Text()) {
			d.inserisci(w)
		}
	}
}

func (d *dizionario) inserisci(w string) {
	isSchema, iniziale := checkIniziale(w)
	if d.checkDuplicati(w, iniziale) {
		return
	}

	seq := tipoSequenza{sequenza: w, isWord: !isSchema}
	(*d.sequenze)[iniziale] = append((*d.sequenze)[iniziale], seq)

	if !isSchema {
		lenParola := len(w)
		(*d.indice)[lenParola] = append((*d.indice)[lenParola], w)
	}
}

func (d *dizionario) elimina(w string) {
	isSchema, iniziale := checkIniziale(w)

	lista := (*d.sequenze)[iniziale]
	for i, c := range lista {
		if c.sequenza == w {
			lista[i] = lista[len(lista)-1]
			(*d.sequenze)[iniziale] = lista[:len(lista)-1]
			break
		}
	}

	if !isSchema {
		lenParola := len(w)
		slice := (*d.indice)[lenParola]
		for i, s := range slice {
			if s == w {
				slice[i] = slice[len(slice)-1]
				(*d.indice)[lenParola] = slice[:len(slice)-1]
				break
			}
		}
		if len((*d.indice)[lenParola]) == 0 {
			delete(*d.indice, lenParola)
		}
	}
}

func (d *dizionario) stampa_parole() {
	Println("[")
	for _, lista := range *d.sequenze {
		for _, seq := range lista {
			if seq.isWord {
				Println(seq.sequenza)
			}
		}
	}
	Println("]")
}

func (d *dizionario) stampa_schemi() {
	Println("[")
	for _, lista := range *d.sequenze {
		for _, seq := range lista {
			if !seq.isWord {
				Println(seq.sequenza)
			}
		}
	}
	Println("]")
}

func compatibile(schema, parola string) bool {
	if len(schema) != len(parola) {
		return false
	}
	mappa := make(map[rune]rune)
	for i, r := range schema {
		corrente := rune(parola[i])
		if unicode.IsUpper(r) {
			if assegnata, ok := mappa[r]; ok {
				if assegnata != corrente {
					return false
				}
			} else {
				mappa[r] = corrente
			}
		} else if r != corrente {
			return false
		}
	}
	return true
}

func (d *dizionario) ricerca(S string) {
	Printf("%s:[\n", S)
	l := len(S)
	for _, parola := range (*d.indice)[l] {
		if compatibile(S, parola) {
			Println(parola)
		}
	}
	Println("]")
}

func distanza(x, y string) int {
	m, n := len(x), len(y)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	maxDistance := m + n
	d := make([][]int, m+2)
	for i := range d {
		d[i] = make([]int, n+2)
	}
	d[0][0] = maxDistance
	for i := 0; i <= m; i++ {
		d[i+1][0] = maxDistance
		d[i+1][1] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j+1] = maxDistance
		d[1][j+1] = j
	}

	da := make(map[rune]int)

	for i := 1; i <= m; i++ {
		db := 0
		for j := 1; j <= n; j++ {
			i1 := da[rune(y[j-1])]
			j1 := db
			cost := 0
			if x[i-1] == y[j-1] {
				cost = 0
				db = j
			} else {
				cost = 1
			}

			d[i+1][j+1] = min(
				d[i][j]+cost,
				d[i+1][j]+1,
				d[i][j+1]+1,
			)

			transCost := d[i1][j1] + (i - i1 - 1) + 1 + (j - j1 - 1)
			if transCost < d[i+1][j+1] {
				d[i+1][j+1] = transCost
			}
		}
		da[rune(x[i-1])] = i
	}
	return d[m+1][n+1]
}

func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= c {
		return b
	}
	return c
}

func (d *dizionario) esisteParola(w string) bool {
	lenParola := len(w)
	for _, parola := range (*d.indice)[lenParola] {
		if parola == w {
			return true
		}
	}
	return false
}

func (d *dizionario) paroleConDistanza1(w string) []string {
	vicine := make([]string, 0)
	lenParola := len(w)
	for _, lenCand := range []int{lenParola - 1, lenParola, lenParola + 1} {
		if lenCand < 0 {
			continue
		}
		for _, candidata := range (*d.indice)[lenCand] {
			if distanza(w, candidata) == 1 {
				vicine = append(vicine, candidata)
			}
		}
	}
	return vicine
}

func (d *dizionario) catena(x, y string) {
	if !d.esisteParola(x) || !d.esisteParola(y) {
		Println("non esiste")
		return
	}
	if x == y {
		Println("(")
		Printf("%s\n", x)
		Println(")")
		return
	}

	queue := []string{x}
	visited := map[string]bool{x: true}
	padre := map[string]string{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, vicina := range d.paroleConDistanza1(current) {
			if visited[vicina] {
				continue
			}
			visited[vicina] = true
			padre[vicina] = current
			if vicina == y {
				d.stampaPercorso(x, y, padre)
				return
			}
			queue = append(queue, vicina)
		}
	}

	Println("non esiste")
}

func (d *dizionario) stampaPercorso(x, y string, padre map[string]string) {
	var percorso []string
	for at := y; at != ""; at = padre[at] {
		percorso = append([]string{at}, percorso...)
		if at == x {
			break
		}
	}
	Println("(")
	for _, parola := range percorso {
		Println(parola)
	}
	Println(")")
}

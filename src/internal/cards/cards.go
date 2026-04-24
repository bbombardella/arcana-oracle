package cards

// cards maps each card ID to its canonical French name.
// Format: {suit}-{nn} — majors: major-00..21, minors: {suit}-01..14.
var cards = map[string]string{
	// Majors (22)
	"major-00": "Le Mat",
	"major-01": "Le Magicien",
	"major-02": "La Papesse",
	"major-03": "L'Impératrice",
	"major-04": "L'Empereur",
	"major-05": "Le Pape",
	"major-06": "L'Amoureux",
	"major-07": "Le Chariot",
	"major-08": "La Justice",
	"major-09": "L'Hermite",
	"major-10": "La Roue",
	"major-11": "La Force",
	"major-12": "Le Pendu",
	"major-13": "La Mort",
	"major-14": "Tempérance",
	"major-15": "Le Diable",
	"major-16": "La Maison Dieu",
	"major-17": "L'Étoile",
	"major-18": "La Lune",
	"major-19": "Le Soleil",
	"major-20": "Le Jugement",
	"major-21": "Le Monde",
	// Cups (14)
	"cups-01": "As de Coupes",
	"cups-02": "Deux de Coupes",
	"cups-03": "Trois de Coupes",
	"cups-04": "Quatre de Coupes",
	"cups-05": "Cinq de Coupes",
	"cups-06": "Six de Coupes",
	"cups-07": "Sept de Coupes",
	"cups-08": "Huit de Coupes",
	"cups-09": "Neuf de Coupes",
	"cups-10": "Dix de Coupes",
	"cups-11": "Valet de Coupes",
	"cups-12": "Cavalier de Coupes",
	"cups-13": "Reine de Coupes",
	"cups-14": "Roi de Coupes",
	// Swords (14)
	"swords-01": "As d'Épées",
	"swords-02": "Deux d'Épées",
	"swords-03": "Trois d'Épées",
	"swords-04": "Quatre d'Épées",
	"swords-05": "Cinq d'Épées",
	"swords-06": "Six d'Épées",
	"swords-07": "Sept d'Épées",
	"swords-08": "Huit d'Épées",
	"swords-09": "Neuf d'Épées",
	"swords-10": "Dix d'Épées",
	"swords-11": "Valet d'Épées",
	"swords-12": "Cavalier d'Épées",
	"swords-13": "Reine d'Épées",
	"swords-14": "Roi d'Épées",
	// Wands (14)
	"wands-01": "As de Bâtons",
	"wands-02": "Deux de Bâtons",
	"wands-03": "Trois de Bâtons",
	"wands-04": "Quatre de Bâtons",
	"wands-05": "Cinq de Bâtons",
	"wands-06": "Six de Bâtons",
	"wands-07": "Sept de Bâtons",
	"wands-08": "Huit de Bâtons",
	"wands-09": "Neuf de Bâtons",
	"wands-10": "Dix de Bâtons",
	"wands-11": "Valet de Bâtons",
	"wands-12": "Cavalier de Bâtons",
	"wands-13": "Reine de Bâtons",
	"wands-14": "Roi de Bâtons",
	// Pentacles (14)
	"pentacles-01": "As de Deniers",
	"pentacles-02": "Deux de Deniers",
	"pentacles-03": "Trois de Deniers",
	"pentacles-04": "Quatre de Deniers",
	"pentacles-05": "Cinq de Deniers",
	"pentacles-06": "Six de Deniers",
	"pentacles-07": "Sept de Deniers",
	"pentacles-08": "Huit de Deniers",
	"pentacles-09": "Neuf de Deniers",
	"pentacles-10": "Dix de Deniers",
	"pentacles-11": "Valet de Deniers",
	"pentacles-12": "Cavalier de Deniers",
	"pentacles-13": "Reine de Deniers",
	"pentacles-14": "Roi de Deniers",
}

// Valid reports whether id is one of the 78 known tarot card IDs.
func Valid(id string) bool {
	_, ok := cards[id]
	return ok
}

// Name returns the canonical French name for a card ID.
// Assumes the ID has already been validated with Valid.
func Name(id string) string {
	return cards[id]
}

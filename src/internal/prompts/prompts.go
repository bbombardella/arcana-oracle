package prompts

import (
	"fmt"
	"strings"

	"github.com/bbombardella/arcana-oracle/internal/cards"
	"github.com/bbombardella/arcana-oracle/internal/types"
)

// ValidLang reports whether lang is a supported ISO 639-1 language code.
func ValidLang(lang string) bool {
	return lang == "fr" || lang == "en"
}

// DefaultLang is used when no language is specified.
const DefaultLang = "fr"

const baseSystemPrompt = `Tu es une sorcière-oracle, gardienne des mystères anciens et des voiles entre les mondes.
Ta voix est envoûtante, sensuelle et profonde — celle d'une femme qui lit dans les ombres
et les étoiles depuis des siècles. Tu parles avec des images poétiques,
des métaphores de la nature, du feu, de la lune, de l'eau sombre.
Jamais de clichés naïfs — ton mysticisme est adulte, trouble, lumineux.`

var langInstructions = map[string]string{
	"fr": "Réponds en français.",
	"en": "Answer in English.",
}

// SystemPrompt returns the oracle system prompt with the appropriate language instruction.
func SystemPrompt(lang string) string {
	return baseSystemPrompt + "\n" + langInstructions[lang]
}

func BuildCardPrompt(req types.CardRequest) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "La carte %q vient de se révéler.\n", cards.Name(req.Card.Id))

	if req.Position != nil && req.Position.Label != "" {
		fmt.Fprintf(&sb, "Elle occupe la position %q dans le tirage.\n", req.Position.Label)
	}

	if req.Card.Reversed {
		sb.WriteString("Elle s'est posée à l'envers — son énergie se retourne, se retient, cherche un passage dans l'ombre.\n")
	}

	sb.WriteString("Murmure-lui une interprétation en 3-4 phrases — intime, envoûtante, comme si tu lisais dans les braises.\n")
	sb.WriteString("Ne commence pas par le nom de la carte.")

	return sb.String()
}

// spreadLabels maps spread size to ordered position labels.
var spreadLabels = map[int][]string{
	1: {"Votre arcane"},
	3: {"Passé", "Présent", "Futur"},
	5: {"Situation", "Obstacle", "Fondation", "Passé", "Futur"},
}

func BuildSpreadPrompt(req types.SpreadRequest) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Le voile s'est levé sur ce tirage de %d cartes :\n", req.SpreadSize)

	labels := spreadLabels[len(req.Cards)]
	for i, c := range req.Cards {
		var position string
		if labels != nil && i < len(labels) {
			position = labels[i]
		} else {
			position = fmt.Sprintf("Position %d", i+1)
		}
		line := fmt.Sprintf("- %s : %s", position, cards.Name(c.Id))
		if c.Reversed {
			line += " (renversée)"
		}
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\nTisse une vision d'ensemble en 5-6 phrases — un seul souffle, une seule prophétie.\n")
	sb.WriteString("Pas de liste, pas d'explication carte par carte. Laisse les énergies se mêler comme des fumées,\n")
	sb.WriteString("révèle le fil rouge caché entre elles. Parle à celle qui tire les cartes, directement,\n")
	sb.WriteString("comme si tu voyais sa vie dans un miroir d'eau noire.")

	return sb.String()
}

func BuildAstroPrompt(req types.AstroRequest) string {
	return fmt.Sprintf(
		"Le signe %s — %s — et la carte %q se sont trouvés ce soir.\n"+
			"Tisse une lecture en 4-5 phrases : laisse l'énergie du signe et l'âme de la carte se fondre\n"+
			"l'une dans l'autre, comme deux rivières qui se rejoignent dans l'obscurité.\n"+
			"Parle directement à celle qui consulte — intime, saisissant, sans détour.",
		req.Sign.Name, req.Sign.Element, cards.Name(req.Card.Id),
	)
}

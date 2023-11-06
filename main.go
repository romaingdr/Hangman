package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	lettresTried  []string
	lettresBonnes []string
	Red           = color.New(color.FgRed)
	Blue          = color.New(color.FgBlue)
	Green         = color.New(color.FgGreen)
	Yellow        = color.New(color.FgYellow)
	Cyan          = color.New(color.FgCyan)
	vie           = 10
)

func main() {
	ClearConsole()
	vie = 10
	lettresTried = []string{}
	lettresBonnes = []string{}
	mot := selectMot()
	fmt.Println(mot)

	for !(motTrouve(mot, lettresTried)) && vie > 0 {
		ClearConsole()
		lettre := playerRound(mot)
		if lettre != "&lostTried" {
			lettresTried = append(lettresTried, lettre)
			if !(lettreDansMot(lettre, mot)) {
				vie--
			}
		} else {
			vie--
		}
	}
	if vie > 0 {
		ClearConsole()
		Green.Println("Vous avez gagné !")
		Cyan.Println("Le mot était : ", mot)
		restart()
	} else {
		ClearConsole()
		affichePendu(vie)
		Red.Println("Vous avez perdu !")
		Cyan.Println("Le mot était : ", mot)
		restart()
	}
}

func selectMot() string {
	var choix int
	var fileName string
	validFile := false
	fmt.Println("Sélectionnez votre liste :")
	fmt.Println("------------------------------------------")
	Cyan.Println("[0] Importer une liste")
	fmt.Println("[1] Mots français")
	fmt.Println("[2] Plus grandes marques")
	fmt.Println("[3] Prénoms les plus courants")
	fmt.Println("[4] Pays")
	fmt.Println("------------------------------------------")
	fmt.Print(">> ")
	fmt.Scanln(&choix)

	rand.Seed(time.Now().Unix())

	switch choix {
	case 0:
		ClearConsole()
		for !validFile {
			fmt.Println("Veuillez saisir le nom du fichier de votre liste")
			fmt.Print(">> ")
			fmt.Scanln(&fileName)

			file, err := os.Open(fileName)
			if err != nil {
				ClearConsole()
				Red.Println("Erreur : le fichier '" + fileName + "' est introuvable.")
			} else {
				file.Close()
				validFile = true
			}
		}
	case 1:
		fileName = "liste-mots/mots.txt"
	case 2:
		fileName = "liste-mots/marques.txt"
	case 3:
		fileName = "liste-mots/prenoms.txt"
	case 4:
		fileName = "liste-mots/pays.txt"
	}

	file, _ := os.Open(fileName)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var mots []string

	for scanner.Scan() {
		mots = append(mots, scanner.Text())
	}

	motAleatoire := mots[rand.Intn(len(mots))]

	return motAleatoire
}

func masquerMot(mot string, lettreTried []string) string {
	motMasque := ""

	for _, lettre := range mot {
		lettreStr := string(lettre)
		if contientLettre(lettreStr, lettreTried) {
			motMasque += lettreStr
			motMasque += " "
		} else {
			motMasque += "_"
			motMasque += " "
		}
	}

	return motMasque
}

func restart() {
	var choix int
	fmt.Println("----------------")
	fmt.Println("[1] Recommencer")
	fmt.Println("[Autre] Quitter")
	fmt.Print(">> ")
	fmt.Scanln(&choix)
	if choix == 1 {
		main()
	} else {
		ClearConsole()
		Red.Println("A bientôt !")
	}
}

func contientLettre(lettre string, lettreTried []string) bool {
	for _, l := range lettreTried {
		if lettre == l {
			return true
		}
	}
	return false
}

func playerRound(mot string) string {
	var lettre string
	longueurMot := len(mot)

	if vie < 10 {
		affichePendu(vie)
	}

	for {
		Cyan.Println("Le mot contient ", longueurMot, " lettres.")
		fmt.Println()
		motMasque := masquerMot(mot, lettresTried)
		fmt.Println(motMasque)
		if len(lettresTried) > 0 {
			fmt.Println()
			fmt.Print("Lettres utilisées : ")
			for i, lettre := range lettresTried {
				if i == 0 {
					if lettreDansMot(lettre, mot) {
						Green.Print(lettre)
					} else {
						Red.Print(lettre)
					}
				} else {
					if lettreDansMot(lettre, mot) {
						Green.Print("; " + lettre)
					} else {
						Red.Print("; " + lettre)
					}
				}
			}
		}
		fmt.Println()

		fmt.Print("Veuillez saisir une lettre : ")
		fmt.Scanln(&lettre)

		if len(lettre) == 1 && isLettreMinuscule(lettre) && !(lettreDansListe(lettre, lettresTried)) {
			break
		} else if lettreDansListe(lettre, lettresTried) {
			ClearConsole()
			if vie < 10 {
				affichePendu(vie)
			}
			Red.Println("Vous avez déjà proposé cette lettre.")
		} else {
			if len(lettre) > 1 {
				if lettre == mot {
					ClearConsole()
					for i := 0; i < len(mot); i++ {
						lettresTried = append(lettresTried, string(mot[i]))
					}
					return string(mot[0])
				} else {
					ClearConsole()
					if vie < 10 {
						affichePendu(vie)
					}
					return "&lostTried"
				}
			} else {
				ClearConsole()
				if vie < 10 {
					affichePendu(vie)
				}
				Red.Println("Veuillez saisir une seule lettre minuscule")
			}
		}
	}
	return lettre
}

func isLettreMinuscule(s string) bool {
	return len(s) == 1 && 'a' <= s[0] && s[0] <= 'z'
}

func ClearConsole() {
	const clearScreen = "\033[H\033[2J"
	fmt.Print(clearScreen)
}

func motTrouve(mot string, lettres []string) bool {
	// Convertit le mot en minuscules pour éviter la casse
	mot = strings.ToLower(mot)

	// Crée une carte pour stocker les lettres de la liste
	lettresPresentes := make(map[string]bool)
	for _, lettre := range lettres {
		lettresPresentes[lettre] = true
	}

	// Parcourt chaque lettre du mot
	for _, lettre := range mot {
		// Convertit la lettre en minuscules pour éviter la casse
		lettreStr := strings.ToLower(string(lettre))
		// Vérifie si la lettre est présente dans la liste
		if !lettresPresentes[lettreStr] {
			return false
		}
	}

	// Si toutes les lettres du mot sont présentes dans la liste, renvoie true
	return true
}

func lettreDansMot(lettre string, mot string) bool {
	mot = strings.ToLower(mot)
	lettre = strings.ToLower(lettre)

	for _, l := range mot {
		if string(l) == lettre {
			return true
		}
	}
	return false
}

func lettreDansListe(lettre string, liste []string) bool {
	// Convertit la lettre en minuscules pour éviter la casse
	lettre = strings.ToLower(lettre)

	// Parcourt chaque lettre dans la liste pour vérifier si la lettre donnée est présente
	for _, l := range liste {
		if strings.ToLower(l) == lettre {
			return true
		}
	}

	// Si la lettre n'est pas trouvée dans la liste, renvoie false
	return false
}

func affichePendu(vie int) {
	switch vie {
	case 9:
		fmt.Println("\n=========")
	case 8:
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 7:
		fmt.Println("  +---+  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 6:
		fmt.Println("  +---+  ")
		fmt.Println("  |   |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 5:
		fmt.Println("  +---+  ")
		fmt.Println("  |   |  ")
		fmt.Println("  O   |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 4:
		fmt.Println("  +---+  ")
		fmt.Println("  |   |  ")
		fmt.Println("  O   |  ")
		fmt.Println("  |   |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 3:
		fmt.Println("  +---+  ")
		fmt.Println("  |   |  ")
		fmt.Println("  O   |  ")
		fmt.Println(" /|   |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 2:
		fmt.Println("  +---+  ")
		fmt.Println("  |   |  ")
		fmt.Println("  O   |  ")
		fmt.Println(" /|\\  |  ")
		fmt.Println("      |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 1:
		fmt.Println("  +---+  ")
		fmt.Println("  |   |  ")
		fmt.Println("  O   |  ")
		fmt.Println(" /|\\  |  ")
		fmt.Println(" /    |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")
	case 0:
		fmt.Println("  +---+  ")
		fmt.Println("  |   |  ")
		fmt.Println("  O   |  ")
		fmt.Println(" /|\\  |  ")
		fmt.Println(" / \\  |  ")
		fmt.Println("      |  ")
		fmt.Println("=========")

	}
}

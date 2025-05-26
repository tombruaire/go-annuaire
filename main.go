package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Contact représenté dans l'annuaire
type Contact struct {
	Nom     string `json:"nom"`
	Prenom  string `json:"prenom"`
	Tel     string `json:"tel"`
}

// Création d'un type Annuaire 
// qui représente la collection de contacts
type Annuaire struct {
	Contacts []Contact `json:"contacts"`
}

const fichierAnnuaire = "annuaire.json"

// Création d'une fonction qui charge l'annuaire depuis le fichier JSON
func chargerAnnuaire() (*Annuaire, error) {
	annuaire := &Annuaire{Contacts: make([]Contact, 0)}
	
	// Vérification si le fichier existe
	if _, err := os.Stat(fichierAnnuaire); os.IsNotExist(err) {
		// Si le fichier n'existe pas, création d'un annuaire vide
		return annuaire, nil
	}
	
	// Lecture du fichier
	data, err := os.ReadFile(fichierAnnuaire)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture du fichier: %v", err)
	}
	
	// Désérialisation du JSON
	err = json.Unmarshal(data, annuaire)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors du parsing JSON: %v", err)
	}
	
	// Affichage de l'annuaire ou null
	return annuaire, nil
}

// Création d'une fonction qui sauvegarde l'annuaire dans le fichier JSON
func (a *Annuaire) sauvegarderAnnuaire() error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("Erreur lors de la sérialisation JSON: %v", err)
	}
	
	err = os.WriteFile(fichierAnnuaire, data, 0644)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'écriture du fichier: %v", err)
	}
	
	return nil
}

// Création d'une fonction qui recherche un contact par nom 
// (insensible à la casse)
func (a *Annuaire) rechercherContact(nom string) *Contact {
	nomLower := strings.ToLower(nom)
	for i := range a.Contacts {
		if strings.ToLower(a.Contacts[i].Nom) == nomLower {
			return &a.Contacts[i]
		}
	}
	return nil
}

// Création d'une fonction qui ajoute un nouveau contact
func (a *Annuaire) ajouterContact(nom, prenom, tel string) error {
	// Vérification si le contact existe déjà
	if a.rechercherContact(nom) != nil {
		// Si le contact existe déjà, affichage d'un message d'erreur
		return fmt.Errorf("Un contact avec le nom '%s' existe déjà", nom)
	}
	
	// Validation des données
	if nom == "" {
		return fmt.Errorf("Le nom ne peut pas être vide !")
	}
	if tel == "" {
		return fmt.Errorf("Le numéro de téléphone ne peut pas être vide !")
	}
	
	// Ajout d'un contact
	contact := Contact{
		Nom:    strings.TrimSpace(nom),
		Prenom: strings.TrimSpace(prenom),
		Tel:    strings.TrimSpace(tel),
	}
	
	a.Contacts = append(a.Contacts, contact)
	return nil
}

// Création d'une fonction qui supprime un contact par nom
func (a *Annuaire) supprimerContact(nom string) error {
	nomLower := strings.ToLower(nom)
	for i, contact := range a.Contacts {
		if strings.ToLower(contact.Nom) == nomLower {
			// Suppression d'un élément du slice
			a.Contacts = append(a.Contacts[:i], a.Contacts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Aucun contact trouvé avec le nom '%s'", nom)
}

// Création d'une fonction qui modifie un contact existant
func (a *Annuaire) modifierContact(nom, nouveauPrenom, nouveauTel string) error {
	contact := a.rechercherContact(nom)
	if contact == nil {
		return fmt.Errorf("Aucun contact trouvé avec le nom '%s'", nom)
	}
	
	// Modification des champs si fournis
	if nouveauPrenom != "" {
		contact.Prenom = strings.TrimSpace(nouveauPrenom)
	}
	if nouveauTel != "" {
		contact.Tel = strings.TrimSpace(nouveauTel)
	}
	
	return nil
}

// Création d'une fonction qui affiche tous les contacts
func (a *Annuaire) listerContacts() {
	if len(a.Contacts) == 0 {
		fmt.Println("Aucun contact dans l'annuaire.")
		return
	}
	
	fmt.Printf("=== Annuaire (%d contact(s)) ===\n", len(a.Contacts))
	for i, contact := range a.Contacts {
		fmt.Printf("%d. %s %s - %s\n", i+1, contact.Nom, contact.Prenom, contact.Tel)
	}
}

func main() {
	// Définition des flags
	var (
		action  = flag.String("action", "", "Action à effectuer (ajouter, lister, rechercher, supprimer, modifier)")
		nom     = flag.String("nom", "", "Nom du contact")
		prenom  = flag.String("prenom", "", "Prénom du contact")
		tel     = flag.String("tel", "", "Numéro de téléphone")
	)
	
	flag.Parse()
	
	// Vérification qu'une action est spécifiée
	if *action == "" {
		fmt.Println("Erreur: Vous devez spécifier une action avec --action")
		fmt.Println("\nActions disponibles:")
		fmt.Println("  ajouter    : Ajouter un nouveau contact")
		fmt.Println("  lister     : Lister tous les contacts")
		fmt.Println("  rechercher : Rechercher un contact par nom")
		fmt.Println("  supprimer  : Supprimer un contact")
		fmt.Println("  modifier   : Modifier un contact existant")
		fmt.Println("\nExemples:")
		fmt.Println("  go run main.go --action ajouter --nom \"Dupont\" --prenom \"Jean\" --tel \"0123456789\"")
		fmt.Println("  go run main.go --action lister")
		fmt.Println("  go run main.go --action rechercher --nom \"Dupont\"")
		fmt.Println("  go run main.go --action supprimer --nom \"Dupont\"")
		fmt.Println("  go run main.go --action modifier --nom \"Dupont\" --prenom \"Pierre\" --tel \"0987654321\"")
		os.Exit(1)
	}
	
	// Chargement de l'annuaire
	annuaire, err := chargerAnnuaire()
	if err != nil {
		fmt.Printf("Erreur lors du chargement de l'annuaire: %v\n", err)
		os.Exit(1)
	}
	
	// Traitement de l'action demandée
	switch strings.ToLower(*action) {
	case "ajouter":
		if *nom == "" || *tel == "" {
			fmt.Println("Erreur: Les paramètres --nom et --tel sont obligatoires pour ajouter un contact")
			os.Exit(1)
		}
		
		err := annuaire.ajouterContact(*nom, *prenom, *tel)
		if err != nil {
			fmt.Printf("Erreur lors de l'ajout: %v\n", err)
			os.Exit(1)
		}
		
		err = annuaire.sauvegarderAnnuaire()
		if err != nil {
			fmt.Printf("Erreur lors de la sauvegarde: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Contact ajouté avec succès: %s %s - %s\n", *nom, *prenom, *tel)
		
	case "lister":
		annuaire.listerContacts()
		
	case "rechercher":
		if *nom == "" {
			fmt.Println("Erreur: Le paramètre --nom est obligatoire pour rechercher un contact")
			os.Exit(1)
		}
		
		contact := annuaire.rechercherContact(*nom)
		if contact == nil {
			fmt.Printf("Aucun contact trouvé avec le nom '%s'\n", *nom)
		} else {
			fmt.Printf("Contact trouvé: %s %s - %s\n", contact.Nom, contact.Prenom, contact.Tel)
		}
		
	case "supprimer":
		if *nom == "" {
			fmt.Println("Erreur: Le paramètre --nom est obligatoire pour supprimer un contact")
			os.Exit(1)
		}
		
		err := annuaire.supprimerContact(*nom)
		if err != nil {
			fmt.Printf("Erreur lors de la suppression: %v\n", err)
			os.Exit(1)
		}
		
		err = annuaire.sauvegarderAnnuaire()
		if err != nil {
			fmt.Printf("Erreur lors de la sauvegarde: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Contact '%s' supprimé avec succès\n", *nom)
		
	case "modifier":
		if *nom == "" {
			fmt.Println("Erreur: Le paramètre --nom est obligatoire pour modifier un contact")
			os.Exit(1)
		}
		
		if *prenom == "" && *tel == "" {
			fmt.Println("Erreur: Au moins un des paramètres --prenom ou --tel doit être fourni pour la modification")
			os.Exit(1)
		}
		
		err := annuaire.modifierContact(*nom, *prenom, *tel)
		if err != nil {
			fmt.Printf("Erreur lors de la modification: %v\n", err)
			os.Exit(1)
		}
		
		err = annuaire.sauvegarderAnnuaire()
		if err != nil {
			fmt.Printf("Erreur lors de la sauvegarde: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Contact '%s' modifié avec succès\n", *nom)
		
	default:
		fmt.Printf("Action inconnue: %s\n", *action)
		fmt.Println("Actions disponibles: ajouter, lister, rechercher, supprimer, modifier")
		os.Exit(1)
	}
}
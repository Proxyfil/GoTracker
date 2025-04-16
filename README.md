# GoTracker

**Réalisé par Pierre-Louis et Margarita**

GoTracker est une application de suivi nutritionnel et de gestion des repas développée en Go. Elle permet aux utilisateurs de suivre leur consommation alimentaire, de calculer des indicateurs comme l'IMC, et de gérer des repas ou des journées prédéfinies.

## Fonctionnalités

- **Recherche d'aliments** : Recherche d'aliments par nom, marque ou catégorie via l'API FoodData Central.
- **Détails nutritionnels** : Affichage des nutriments pour une quantité donnée.
- **Gestion des utilisateurs** : Inscription, connexion et mise à jour des informations utilisateur.
- **Suivi des indicateurs** : Historique de l'IMC, du poids et du pourcentage de graisse corporelle.
- **Gestion des repas et journées** : Création, ajout et gestion de repas et de journées prédéfinies.
- **Historique alimentaire** : Suivi des aliments consommés.

## Prérequis

- **Go** : Version 1.24 ou supérieure.
- **PostgreSQL** : Version 15 ou supérieure.
- **Docker** (optionnel) : Pour exécuter PostgreSQL dans un conteneur.

## Installation

### 1. Cloner le dépôt

```bash
git clone https://github.com/Proxyfil/GoTracker.git
cd GoTracker
```

### 2. Configurer la base de données

#### Option 1 : Utiliser Docker
Lancez un conteneur PostgreSQL avec Docker :

```bash
docker-compose up -d
```

#### Option 2 : Configurer PostgreSQL manuellement
Assurez-vous que PostgreSQL est installé et configurez les informations de connexion dans le fichier `src/vars/config.json` :

```json
{
  "db_user": "postgres",
  "db_password": "mypassword",
  "db_host": "localhost",
  "db_port": "5433",
  "db_name": "gotracker"
}
```

### 3. Installer les dépendances

Dans le répertoire `src`, exécutez :

```bash
go get github.com/lib/pq
```

### 4. Lancer l'application

Dans le répertoire `src`, exécutez :

```bash
go run .
```

## Utilisation

### Commandes disponibles

Voici une liste des commandes disponibles dans l'application :

- **`help`** : Affiche la liste des commandes disponibles.
- **`register <firstname> <lastname> <age> <weight> <height> <target_weight>`** : Inscrit un nouvel utilisateur.
- **`login <user_id>`** : Connecte un utilisateur existant.
- **`bodyfat`** : Affiche le pourcentage de graisse corporelle de l'utilisateur connecté.
- **`imc`** : Affiche l'IMC de l'utilisateur connecté.
- **`search food <food_name>`** : Recherche un aliment par nom.
- **`search_with_filter <food_name> <dataType>`** : Recherche un aliment avec un filtre (ex. : `Foundation`).
- **`search_by_brand_or_category <food_name> <brandOwner> <foodCategory>`** : Recherche un aliment par marque ou catégorie.
- **`details <food_id>`** : Affiche les détails nutritionnels d'un aliment.
- **`add food <food_id> <quantity>`** : Ajoute un aliment consommé à l'historique.
- **`create meal <meal_name> <meal_type>`** : Crée un nouveau repas.
- **`list meal`** : Liste tous les repas.
- **`exit`** : Quitte l'application.

### Exemple d'utilisation

1. **Inscription d'un utilisateur** :
   ```bash
   register John Doe 30 70 175 65
   ```

2. **Connexion d'un utilisateur** :
   ```bash
   login 1
   ```

3. **Recherche d'un aliment** :
   ```bash
   search food apple
   ```

4. **Afficher les détails nutritionnels** :
   ```bash
   details 454004
   ```

5. **Ajouter un aliment consommé** :
   ```bash
   add food 454004 150
   ```

## Structure du projet

```
GoTracker/
├── docker-compose.yml       # Configuration Docker pour PostgreSQL
├── README.md                # Documentation du projet
├── src/
│   ├── main.go              # Point d'entrée principal
│   ├── cli/                 # Gestion des commandes CLI
│   │   └── cli.go
│   ├── fdcnal/              # Intégration avec l'API FoodData Central
│   │   └── api.go
│   ├── utils/               # Utilitaires (connexion à la base de données, etc.)
│   │   └── db.go
│   ├── structs/             # Structures de données
│   │   └── user.go
│   ├── interface/           # Interfaces utilisateur
│   │   └── user.go
│   └── vars/                # Variables de configuration
│       └── config.json
```

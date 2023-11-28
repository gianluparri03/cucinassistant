![CucinAssistant](application/static/banner.png)

Una piattaforma web (multiutente) creata per gestire:

- Lista della spesa
- Lista delle idee
- Menu della settimana
- Quantità
- Scadenze

## Installazione

Una versione di CucinAssistant è disponibile online per il pubblico (controllare su github per il link).

Se si vuole hostare la propria istanza, CucinAssistant può essere eseguito in due modi:

1. **Docker**

Scaricando l'immagine da github con `docker pull ghcr.io/gianluparri03/cucinassistant`.
Il server rimarrà in ascolto sulla porta `80` del container.


2. **Da sorgente**

- Installando tutte le dipendenze con `pip install -r requirements.txt`
- Creando le variabili d'ambiente `PRODUCTION=1` e `SECRET=StringaACaso`
- Eseguendo `run.py`

Il server rimarrà in ascolto sulla porta `80`.

**In ogni caso**, tutti i dati vengono salvati nella cartella `application/data/` (`/cucinassistant/application/data` sull'immagine Docker)

## Limitazioni

Al momento, non disponendo di un server web, il cambio o il reset della password non può essere fatto in automatico, ma deve
essere effettuato da un amministratore che ha accesso al server (le password rimangono comunque illeggibili, in quanto criptate).

## Crediti

Questo sito web utilizza [sakura](https://github.com/oxalorg/sakura) come framework css.

## Altro

Il progetto è rilasciato con la [licenza MIT](/blob/main/LICENSE).

Qualsiasi pull request è ben accetta; per ogni domanda, dubbio o proposta, potete contattarmi su <a href="mailto:gianluparri03@gmail.com?subject=[CucinAssistant]">gianlucaparri03@gmail.com</a>.

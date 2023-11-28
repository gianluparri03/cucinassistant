![CucinAssistant](banner.png)

Una piattaforma web (multiutente) creata per gestire:

- Lista della spesa
- Lista delle idee
- Menu della settimana
- Quantità
- Scadenze

## Installazione

Una versione di CucinAssistant è disponibile online per il pubblico (controllare su github per il link).

Se si vuole hostare la propria istanza, si può fare in due modi:

1. **Docker**

L'immagine Docker può essere scaricata con

`docker pull ghcr.io/gianluparri03/cucinassistant`

2. **Da sorgente**

Installando tutte le dipendenze (`pip install -r requirements.txt`) ed eseguendo il file `run.py`; il server rimarrà in ascolto
su `0.0.0.0:8080`. Impostando la variabile di ambiente `PRODUCTION=1` si disattiverà la modalità debug di Flask.

## Limitazioni

Al momento, non disponendo di un server web, il cambio o il reset della password non può essere fatto in automatico, ma deve
essere effettuato da un amministratore che ha accesso al server (le password rimangono comunque illeggibili, in quanto criptate).

## Altro

Il progetto è rilasciato con la [licenza MIT](/blob/main/LICENSE).

Qualsiasi pull request è ben accetta.

Per ogni domanda, dubbio o proposta, potete contattarmi su <a href="mailto:gianluparri03@gmail.com?subject=[CucinAssistant]">gianlucaparri03@gmail.com</a>.

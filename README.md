Rassegna Stampa CNI
===================

Si tratta di uno scraper/mailer che scarica la [rassegna stampa del Centro
Studi CNI](http://www.centrostudicni.it/rassegna-stampa) ed invia i PDF in allegato via mail. Il Centro Studi CNI offre 
già un servizio di mailing list per la rassegna stampa, tuttavia bisogna
invece di allegare il PDF nelle email inseriscono un collegamento ad una
pagina web in cui a sua volta c'è il collegamento al PDF. Consiglio caldamente
di usare il loro servizio, tuttavia questa era un'occasione di provare
un po' di Go. :-)

Per compilare questo programma è necessario scaricare tramite "go get"
l'ottima libreria [GoQuery](https://github.com/PuerkitoBio/goquery).

    go get github.com/PuerkitoBio/goquery

Bisogna anche scaricare le meno ottime ma pur sempre necessarie "goutils"
del sottoscritto. ;-)

    go get github.com/bernarpa/goutils

Per far funzionare il programma bisogna creare un file di configurazione (di solito
è chiamato rscni.cfg) ed impostare nella variabile d'ambiente RSCNICFG
il percorso completo di quel file (incluso il file stesso, ad esempio
/home/pippo/.rscni.cfg). Nella directory cfg c'è un esempio di questo file.

Dentro il file basta impostare alcune direttive. In particolare:
   * datadir: la directory dove verranno salvati i file di rscni
   * smtp.*: i parametri di configurazione del server smtp per inviare
             le mail (attualmente si presuppone che il server funzioni
			 tramite TLS).

Inoltre, nella datadir, è necessario predisporre un file chiamato ml.txt
contenente gli indirizzi email a cui verrà inviata la rassegna stampa.
Ciascuna riga di ml.txt è come quella seguente:

    Nome Cognome <indirizzo@email.com>

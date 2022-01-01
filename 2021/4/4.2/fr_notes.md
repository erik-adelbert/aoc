Première remarque ici, utiliser une [`pile`](https://yourbasic.org/golang/implement-stack/) n'est pas nécessaire pour résoudre le problème: il suffirait de mémoriser la dernière carte gagnante.

J'ai préféré utiliser une `pile` parce que le résultat a une portée plus générale: dans la pile il y a *l'historique complet* de la partie et il devient possible de répondre à n'importe quelle question sur l'ordre des cartes gagnantes par rapport au tirage. p.ex. je sais qu'à tout moment de la partie, la dernière carte gagnante est sur le sommet de la `pile`.

À mesure que les cartes sont gagnantes, il faut aussi les enlever de notre `deck`. Mais, en `Go` *retailler* une `slice` *pendant qu'on itère dessus* est *indéfini* (en pratique ça ne marche pas): on ne peut pas supprimer les cartes du jeu à mesure q'elles sont gagnantes.  
On peut contourner cette limitation 1) si on ne modifie pas la taille de la `slice` pendant l'itération et 2) si on écrit au début, *avant* (ou *sur*) le pointeur courant

C'est ce qui se passe entre les lignes 108~119: pour chaque numéro tiré, soit la carte est gagnante et elle va dans la `stack`, soit elle reste dans notre jeu et on la remet au début du `deck`. Lorsqu'on a fini une passe sur le `deck`, on peut retailler sa `slice` avec le nombre des cartes restantes.  
Ça fonctionne parce qu'au lieu de *supprimer* les cartes gagnantes, on *conserve* les autres; on voit que `i` de la ligne 109 est toujours plus petit que l'indice implicite de la ligne 110; l'ordre relatif des cartes est aussi préservé (c'est une bonne propriété qui vient gratuitement).
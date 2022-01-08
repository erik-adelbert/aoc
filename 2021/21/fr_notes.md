Le problème décrit un jeu de plateau à deux où on joue chacun à son tour. Dans part2, il y a toujours 27 coups possibles mais la partie est limitée par la faiblesse du score gagnant (21). 

L'algorithme [récursif](https://en.wikipedia.org/wiki/Recursion_(computer_science)) qui résoud ce problème est le premier qu'on étudie en *théorie des jeux* mais plus généralement il concerne la prise de décision: il s'agit de [*minimax*](https://en.wikipedia.org/wiki/Minimax). Il exploite utilement l'idée que dans un jeu à deux, une partie c'est un premier coup du joueur au trait, suivi d'une partie où c'est l'autre joueur qui commence, jusqu'à la victoire.

Ici, il fonctionne bien parce qu'il n'y a pas d'information secrète (pas de dé ou de position cachée), le jeu est à *information complète*. Comme la victoire d'un joueur entraîne la défaite de l'autre, le jeu est *à somme nulle*. Comme les coups dépendent d'une petite combinatoire et qu'on peut tous les générer: on peut le *résoudre totalement* càd calculer tous les jeux possibles.

Un `état du jeu` est un vecteur `(c1, s1, c2, s2)` avec `c1` et `s1` la position et le score du joueur `p1`. Une `victoire` est un état qui comporte un score `s1` ou `s2` supérieur à 21, il n'y a pas de coup au-delà. Une `partie` est un ensemble d'états reliés par des coups jusqu'à une victoire.  
Pour `résoudre` le jeu, à la manière de `minimax`, on commence avec `(c1, 0, c2, 0)` et on joue successivement toutes les parties qui découlent des coups possibles.  
Pour jouer une partie, à partir de `(c1, s1, c2, s2)`, on joue un coup pour `p1`, on vérifie s'il est gagnant, sinon on met à jour l'état courant `(c1, s'1, c2, s2)` et on joue toutes les sous-parties à partir de `(c2, s2, c1, s'1)` (c'est `p2` qui commence) avant de passer au coup suivant de `p1` et de recommencer.  

Quand on fait ça on construit complètement [*l'arbre du jeu*](https://en.wikipedia.org/wiki/Game_tree), on dit qu'on *résoud totalement* le jeu. Il y a très peu de [jeux](https://en.wikipedia.org/wiki/Hex_(board_game)) qu'on peut résoudre totalement.

C'est le premier algorithme qu'on étudie parce qu'il est lié au [théroème](https://en.wikipedia.org/wiki/Minimax_theorem) qui fonde la théorie des jeux et qu'on doit à [john von neumann](https://en.wikipedia.org/wiki/John_von_Neumann) lui-même.

Le dernier né (et le plus impressionnant) de cette théorie est le programme [`α0`](https://en.wikipedia.org/wiki/AlphaZero). Son niveau excède largement le nôtre; il a créé de la [connaissance](https://deepmind.com/blog/article/alphazero-shedding-new-light-grand-games-chess-shogi-and-go), inconnue auparavant, au sujet du jeu de [go](https://en.wikipedia.org/wiki/Go_(game)).

<div style="text-align:center">
  <img src="https://www.ocf.berkeley.edu/~yosenl/extras/alphabeta/alphabeta.jpg" />
</div>
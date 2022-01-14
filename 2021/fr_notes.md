* * *
## Day 1

Pour cette inauguration d'aoc2021, j'ai écrit et soumis ce programme en moins de 5mn mais ça ne m'a fait entrer que dans le top 3000, wow!  
Il s'agit de compter, dans une suite de nombres, combien de fois on a une augmentation entre deux nombres successifs. Rien à dire sur le fond, ici c'est la vitesse de composition qui compte.  
Je choisis d'utiliser 3 variables (plutôt qu'un tableau p.ex.) parce que *the simpler, the better*: à lisibilité et performance égale, j'essaie de choisir la forme la plus simple. Cela dit, 3 c'est la limite au-delà de laquelle le tableau est mieux.  
  
En revanche, sur la forme, un petit programme comme celui-ci c'est l'occasion de sortir, du premier coup, un programme *bien écrit*. C'est à dire avec une nomenclature et un style choisis pour durer: la construction `for input.scan(){}` est un classique qu'on retrouve à l'identique dans le manuel du Go par [Donovan & Kernighan (D&K)](https://www.gopl.io). C'est effficient parce que ça énonce simplement, clairement et de manière compacte qu'elle est l'intention de programmation: on veut scanner les entrées.  
En ce qui concerne la variable `old`, elle pourrait s'appeler `last`, c'est plus traditionnel. Mais, en pratique, il m'est plus facile de manipuler des symboles d'à peu près la même longueur. Ici, on a une variable "courante" et une variable "précédente" soient: `last` (ou `prev`) et `cur` ou plus court, `old` et `cur` qui aident tout aussi bien à capturer le sens de ce qui est fait au premier coup d'oeil.  
Il y a aussi le minimalisme de `Go`, qui m'oblige à déclarer la constante `MaxInt` de façon [*idiomatique*](https://dgryski.medium.com/idiomatic-go-resources-966535376dba).

* * *
## Day 2

Dans les problèmes de simulation, souvent, on peut faire quelquechose sans avoir à capturer toute l'entrée: ici, p.ex., les lignes d'entrée se composent d'une commande suivi d'un argument (nombre). Les initiales des commandes sont différentes les unes des autres: lire `line[0]` suffit à décoder une commande (`f`, `u`, `d`). Enfin, lorsqu'on découpe `line` sur son espace central, à droite il y a le nombre.

* * *
## Day 3

J'adore le support `utf8` de `Go`, les γράμματα grecques et aussi les programmes *vintage* bien bas niveau (bits): je me suis fait plaisir gratuitement. Sinon, ce programme ne supporte que les nombres 12bits, il pourrait supporter une largeur arbitraire dynamique (pour presque rien) mais j'ai préféré la simplicité: la largeur est constante; elle est facile à éditer si les entrées changent.

Pour la `part2`, l y a une unique fonction `rate()`. Elle exécute une de ses deux branches en fonction de son entrée. Comme le choix est binaire, c'est un booléen qui décide son mode. C'est pour rendre ce booléen plus lisible et faciliter la mise au point que je déclare les deux constantes (o2 <- O2, co2 <- CO2). 

`rate()` mesure (sous la forme d'une `string`) les *most/least popular bits* des `inputs` (`strings` aussi) et retourne le résultat de `strconv.ParseInt()`, de la bibliothèque standard, sur cette mesure.
Au lieu de gérer ou filtrer `err` sur place (dans `rate()`), je la laisse remonter jusqu'à ce que je sois obligé de la gérer: juste avant de la transmettre sur son channel. C'est ce que je préfère comme gestion d'erreur: mes programmes ne gèrent une erreur que quand elle ne peut plus remonter. À ce moment-là, elle a souvent un sens bien défini.

Comme les problèmes sont bien séparés, je demande l'éxecution [concurrente](https://youtu.be/oV9rvDllKEg) de `rate()` dans deux `goroutines`: c'est gratuit!  

Enfin, [`popcnt`](https://en.wikipedia.org/wiki/SSE4#POPCNT_and_LZCNT) est une instruction `CPU` qui compte les bits à un dans un entier. C'est H.S. Warren Jr. qui a popularisé ce nom dans [Hacker's Delight](https://en.wikipedia.org/wiki/Hacker%27s_Delight). En vrai, il s'agit du [*poids de Hamming*](https://en.wikipedia.org/wiki/Hamming_weight) et j'utilise le nom `popcount` ou `popcnt` pour tous les problèmes où il faut dénombrer une population.

* * *
## Day 4

Dans ce problème, le programme doit jouer au loto. C'est à dire tenir le compte des numéros qui ont été tirés sur ses cartes.
Pour modéliser une carte, j'utilise une structure qui groupe:
- 5 compteurs de lignes
- 5 compteurs de colonnes
- 1 `map` des nombres et de leur position sur la carte

p.ex. cette carte (en 3x3 au lieu de 5x5):

<table>
<tbody>
  <tr>
    <td>1</td>
    <td>2</td>
    <td>3</td>
  </tr>
  <tr>
    <td>4</td>
    <td>5</td>
    <td>6</td>
  </tr>
  <tr>
    <td>7</td>
    <td>8</td>
    <td>9</td>
  </tr>
</tbody>
</table>


est modélisée comme ceci:
<table>
<tbody>
  <tr>
    <td></td>
    <td>12</td>
    <td>15</td>
    <td>18</td>
  </tr>
  <tr>
    <td>6</td>
    <td>.</td>
    <td>.</td>
    <td>.</td>
  </tr>
  <tr>
    <td>15</td>
    <td>.</td>
    <td>.</td>
    <td>.</td>
  </tr>
  <tr>
    <td>24</td>
    <td>.</td>
    <td>.</td>
    <td>.</td>
  </tr>
</tbody>
</table>

{  
    1: (1, 1), 2: (1, 2), 3: (1, 3),  
    4: (2, 1), 5: (2, 2), 6: (2, 3),  
    7: (3, 1), 8: (3, 2): 9: (3, 3),  
}  

S'il sort `3`, alors:
- on lit les coordonnées de `3` sur la `map`: `{1, 3}`
- on efface `3` de la `map`
- on met à jour le compteur de la ligne `1`: `6-3`
- on met à jour le compteur de la colonne `3`: `18-3`

<table>
<tbody>
  <tr>
    <td></td>
    <td>12</td>
    <td>15</td>
    <td><b>15</b></td>
  </tr>
  <tr>
    <td><b>3</b></td>
    <td>.</td>
    <td>.</td>
    <td>.</td>
  </tr>
  <tr>
    <td>15</td>
    <td>.</td>
    <td>.</td>
    <td>.</td>
  </tr>
  <tr>
    <td>24</td>
    <td>.</td>
    <td>.</td>
    <td>.</td>
  </tr>
</tbody>
</table>

{  
    1: (1, 1), 2: (1, 2), *~~3: (1, 3)~~*,  
    4: (2, 1), 5: (2, 2), 6: (2, 3),  
    7: (3, 1), 8: (3, 2): 9: (3, 3),  
}  

Si un compteur tombe à zéro, la ligne ou colonne correspondante est gagnante, c'est [`bingo`](https://fr.wikipedia.org/wiki/Loto#Bingo)! Si la somme des compteurs de lignes (ou de colonnes) vaut zéro ou que la `map` est vide, la carte est gagnante.

À chaque `bingo`, pour calculer la somme des nombres restants sur la carte, il suffit de sommer soit les compteurs des lignes, soit ceux des colonnes.

La structure de donnée *simplifie* la résolution en supportant *facilement* et *rapidement* les opérations du tirage de loto. Comme l'entrée ne contient pas de doublon, je ne vide pas la `map` dans la fonction `biff()`.

Utiliser une [`pile`](https://yourbasic.org/golang/implement-stack/) n'est pas nécessaire pour résoudre le problème: il suffirait de mémoriser la première et la dernière carte gagnante.

J'ai préféré en implémenter une parce que le résultat a une portée plus générale: dans la pile il y a *l'historique complet* de la partie et il devient possible de répondre à n'importe quelle question sur l'ordre des cartes gagnantes: p.ex. à tout moment du tirage, la dernière carte gagnante est sur le sommet de la `pile`.

En `Go` retailler une `slice` pendant qu'on itère dessus est *indéfini* (en pratique, ça ne fonctionne pas): on ne peut pas directement supprimer les cartes du jeu à mesure q'elles sont gagnantes.

On peut contourner cette limitation 1) si on ne modifie pas la taille de la `slice` pendant l'itération et 2) si on écrit au début, *avant* ou *sur* le pointeur courant. Ici, entre les lignes 88~97 on voit que pour chaque numéro tiré, si la carte est gagnante et elle va dans la `stack`, sinon elle `retourne` dans notre `deck`. Lorsqu'on a fini une passe, on peut retailler la `slice` du `deck` avec le nombre des cartes restantes.  

Ça fonctionne parce qu'au lieu de *supprimer* les cartes gagnantes, on *conserve* les autres; on voit que `i` de la ligne 88 est toujours plus petit que l'indice implicite de la ligne 89; l'ordre relatif des cartes est aussi préservé (c'est une bonne propriété qui vient gratuitement).

* * *
## Day 5

Si les lignes avaient été autrement qu'à ±45º, il aurait fallu utiliser l'algorithme de [bresenham](https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm)! 

* * *
## Day 6

C'est en tournant l'exemple à la main que j'ai vu qu'il s'agissait d'une rotation vers la gauche avec une addition sur le jour 6. Je trouve la version `python` délicieuse!

* * *
## Day 7

Une des plus nice [contributions](https://www.reddit.com/r/adventofcode/comments/rawxad/2021_day_7_part_2_i_wrote_a_paper_on_todays/?utm_source=share&utm_medium=web2x&context=3) sur reddit, cette année!

Pour l'histoire, pendant la première partie, je me suis trompé et j'ai commencé par prendre la moyenne au lieu du médian. À la lecture de la seconde partie, j'ai reconnu la moyenne dans l'exemple. Pour mon étoile, j'ai utilisé `ceil()` qui marchait sur mes `inputs` mais j'étais pas sûr du tout que se soit pas `floor()` puisque c'était à l'instinct. Et comme `ceil()` ne fonctionnait pas (à peu de chose près) tandis que `floor()` si, sur d'autres `inputs`, j'en ai déduit que c'était `round()` la solution. Quand j'ai lu le papier, le lendemain, ça m'a mis de bonne humeur!

* * *
## Day 9

J’ai traité la part2 avec [ceci](https://www.ocf.berkeley.edu/~fricke/projects/hoshenkopelman/hoshenkopelman.html) que j’avais déjà utilisé dans [un autre projet](https://github.com/erik-adelbert/mcs/blob/master/pkg/chaingame/tag.go). Cet algorithme s’appuie sur un classique union-find ([ici](https://www.cs.princeton.edu/~rs/AlgsDS07/01UnionFind.pdf), le cours de [sedgewick](https://en.wikipedia.org/wiki/Robert_Sedgewick_(computer_scientist))) qui permet de remédier aux problèmes sur certains contours et avec les composants “concaves” engendrés par un `flooding nw`.  
Avec *Hoshen-Kopelman* on obtient la sortie parfaite sans pour autant faire exploser la complexité de l’algo (toutes les parties de l’algo sont linéaires en `n` avec `n` le nombre d’entrées) et surtout sans avoir à gérer une multitude de cas incroyables (impossible à mettre au point). Ma solution du jour pour part2 s'exécute en temps linéaire \o/

* * *
## Day 10

Il s’agit d’un classique qui a même un nom: `bracket matching`. Il y en a dans tous nos `IDE` et traditionnellement on peut le résoudre à l’aide d’une [pile](https://www.geeksforgeeks.org/check-for-balanced-parentheses-in-an-expression/).

Ici, au lieu de pousser le *symbole ouvrant* dans la `stack`, j'envoie le symbole fermant correspondant (ligne 48): ça simplifie le test d'appariement qui arrive plus tard (ligne 50).

* * *
## Day 11

Pur exercice de [programmation dynamique](https://en.wikipedia.org/wiki/Dynamic_programming) et de [*multi-buffering*](https://en.wikipedia.org/wiki/Multiple_buffering).

* * *
## Day 12

Pur exercice de [recherche en profondeur d'abord](https://en.wikipedia.org/wiki/Depth-first_search) sur un graphe. Ça sera un problème de bac série `code`, dans le turfu. 

* * *
## Day 13

Ce problème déguise, sous une forme romancée, deux techniques de base, une de la [*théorie des codes*](https://en.wikipedia.org/wiki/Coding_theory) et une de [*l'informatique graphique*](https://en.wikipedia.org/wiki/Computer_graphics): les points sont donnés sous une forme vectorielle. Après 1) les avoir transformés (*décodés*) conformément à un modèle (*code*) de pliage, il faut 2) les afficher à l'écran.

J'ai traité cet affichage comme un [*Raster Scan Display*](https://www.geeksforgeeks.org/raster-scan-displays/). Pour l'implémenter, j'utilise une [*aabb*](https://en.wikipedia.org/wiki/Minimum_bounding_box#Axis-aligned_minimum_bounding_box) combinée avec un [*framebuffer*](https://en.wikipedia.org/wiki/Framebuffer). Pour obtenir une *image*, je [*rasterise*](https://en.wikipedia.org/wiki/Rasterisation) les points dans ce *buffer*. Grâce à un teammate, pendant une [review](https://en.wikipedia.org/wiki/Code_review), j'ai découvert que le caractère `undefined` `�` de l'`ASCII` étendu est un des plus lumineux. C'est pourquoi je le choisis comme valeur du *pixel allumé*. 

Quand j'utilise des [*abstractions connues*](https://en.wikipedia.org/wiki/Abstraction_(computer_science)), je gagne du temps: je facilite l'écriture, la mise au point et la maintenance du programme; je sais à l'avance ce qu'il faut faire et comment il faut le faire parce que la documentation est abondante. Si je les utilise suffisamment, je finis par les connaître par coeur.

Je trouve ce problème très fun, j'ai bien aimé écrire un programme [efficace](https://en.wikipedia.org/wiki/Algorithmic_efficiency) pour le résoudre.

```
❯ cat input.txt| ./aoc13.2
�    ���  ����   �� ���    �� ���� �  �
�    �  � �       � �  �    � �    �  �
�    �  � ���     � ���     � ���  ����
�    ���  �       � �  �    � �    �  �
�    � �  �    �  � �  � �  � �    �  �
���� �  � �     ��  ���   ��  ���� �  �
```

`Side Note:`  
Quand on résoud ce problème, on trouve que pour un point `p(x,y)` qui bouge pendant une [symétrie](https://en.wikipedia.org/wiki/Reflection_symmetry) d'axe `a` on a:  
$$x_{p_{n+1}} = 2*x_{a} - x_{p_{n}}\ ||\ y_{p_{n+1}} = 2*y_{a} - y_{p_{n}}\\$$

On peut montrer pourquoi, à l’aide des coordonnées homogènes de *moebius*. Dans ce système, ces [matrices](https://en.wikipedia.org/wiki/Transformation_matrix) représentent une *translation* de vecteur `u(x, y)`, une *symétrie* sur l’axe des `x` (pliure horizontale) et une sur celle des `y` (pliure verticale):  

$$\begin{pmatrix}1 & 0 & x_{u}\\0 & 1 & y_{u}\\0 & 0 & 1\\\end{pmatrix}
\begin{pmatrix}1 & 0 & 0\\0 & -1 & 0\\0 & 0 & 1\\\end{pmatrix}
\begin{pmatrix}-1 & 0 & 0\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix}$$ 

Plier la feuille en deux verticalement, sur l’axe `x=a`, ça équivaut à 1) translater les points concernés de `(-a, 0)` (pour les mettre dans un repère centré sur eux), 2) les symétriser verticalement et 3) les translater de `(a, 0)` (pour les remettre dans le repère d’origine). Soit une combinaison des transformations suivantes: $$Tr_{a}(x).S_{y}(x).Tr_{-a}(x)$$ Ça se lit de droite à gauche, on ne peut pas changer l’ordre. Pour transformer un point, on écrit:
$$\begin{pmatrix}1 & 0 & x_{a}\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} . 
\begin{pmatrix}-1 & 0 & 0\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} .
\begin{pmatrix}1 & 0 & -x_{a}\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} .
\begin{pmatrix}x_{p}\\y_{p}\\z_{p}\end{pmatrix}$$

Et on voit apparaître la formule: $$x_{p_{n+1}} = 2*x_{a} - x_{p_{n}}$$  
Pour des problèmes où tous les points sont concernés, ça devient très puissant: on peut souvent combiner (multiplier) toutes les transformations en une seule matrice et ensuite l’appliquer en masse (`O(T+N)`) au lieu de calculer chaque transformation pour chaque point (`O(T*N)`)…

* * *
## Day 15

Lorsqu'on utilise [*l'algorithme de Dijkstra*](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm), on ne peut plus dire grand chose d'autre: ça construit un genre de [*mycelium*](https://en.wikipedia.org/wiki/Mycelium) de données, à mesure que l'algorithme décompose (grignote) le problème. Pour accélerer le processus et guider la recherche (le grignotage), on utilise une [*priority queue*](https://en.wikipedia.org/wiki/Priority_queue). Ici, j'utilise [celle](https://pkg.go.dev/container/heap) de la bibliothèque standard de `Go`.
  
<div style="text-align:center">
  <img src="https://upload.wikimedia.org/wikipedia/commons/2/23/Dijkstras_progress_animation.gif" />
</div>

* * *
## Day 16

Pour résoudre ce problème linéairement, j'utilise l'excellente bibliotèque [`bitstream-go`](https://github.com/bearmini) de `bearmini` que je combine avec `math/big` de la bibliotèque standard. 

Ces bibliothèques implémentent des [interfaces](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go) standardisées comme [`Reader`](https://go.dev/tour/methods/21) et [`Writer`](https://www.grant.pizza/blog/the-beauty-of-io-writer/).

Grâce à ces interfaces, les (fonctions) filtres ont le même *prototype* (interface) que les fonctions bas niveau devant laquelle elle sont posées. Avoir des fonctions avec des prototypes identiques en entrée et en sortie permet de fabriquer des [*pipelines de données*](https://en.wikipedia.org/wiki/Pipeline_(computing)) comme sur la ligne 161.

* * *
## Day 18

Le problème nous présente une variété d'[arbres binaires](https://en.wikipedia.org/wiki/Binary_tree#Internal_nodes) sur lesquelles sont définis des opérations: les *snailfish numbers*.

Ces arbres sont composés de *feuilles* qui contiennent un entier, ces feuilles sont toutes reliées entre elles (et à la *racine* de l'arbre) par des *noeuds internes* qui ne contiennent q'une unique *paires* de *liens*. Cette structure de données est aussi connue pour opérer une *classification* comme dans les [`k-d trees`](https://en.wikipedia.org/wiki/K-d_tree) ou les [`B-trees](https://en.wikipedia.org/wiki/B-tree).

Pour la lecture des entrées, j'injecte suffisamment d'espaces dans l'entrée pour pouvoir capturer *chaque* symbole séparemment en *une seule fois*.

La fonction `newPair()` a des *arguments variables*.

L'opération `explode()` tire utilement partie de la forme [*aplatie*](https://www.geeksforgeeks.org/flatten-a-binary-tree-into-linked-list/) de l'arbre pour mettre à jour des feuilles *adjacentes*; il n'y a que sous cette forme que l'information de *voisinage* est disponible.

Dans l'opération `reduce()`, j'utilise des [`drapeaux binaires`](https://en.wikipedia.org/wiki/Mask_(computing)) pour *synchroniser* le *worflow* (`done` ligne 174~184).

Enfin, j'ai profilé la version de base et découvert que l'essentiel du *runtime* de part2 consistait à attendre pour fournir des *hints* au kernel au sujet de l'utilisation *mémoire*. J'ai pris la décision de rendre la totalité de part2 concurrente: Je lance un producteur qui lance des *sous-producteurs* et en même temps, je lance des *consumers* qui calculent les magnitudes. Dans la routine `main`, je collecte et filtre les résultats. Le résultat net de cette transformation est de ramener le *runtime* de 651ms à 245ms!

`MàJ` Tout ce qui est dit plus haut est vrai mais trop compliqué: en repensant ce problème depuis le début, j'ai trouvé qu'en stockant les valeurs et les profondeurs des *snail numbers* la structure de donnée rendait toutes les opérations plus faciles à l'exception du calcul de *magnitude* qui a néanmoins une difficulté acceptable. Au total j'aurais passé plus de 24h cumulés sur ce sujet. Mais c'est lui qui me fait passer sous la barre de la seconde sur mon mb air m1 \o/

* * *
## Day 19

Le sujet du problème est effrayant au premier abord.

J'ai voulu une solution mécanique. J'ai *précalculé* les rotations et je les ai modelisées en deux parties: les *signes* des rotations et l'*ordre des axes*. Comme d'habitude, quand je fais ça, j'ai bien saigné des neurones et j'ai eu recours à une boîte d'allumettes et un stylo pour la mise au point.

* * *
## Day 21

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

* * *
## Day 22

J'ai mis des heures à choisir une structure de données qui fonctionne et des heures à mettre au point ce programme. C'est une implémentation *minimaliste* des [*k-d trees*](https://en.wikipedia.org/wiki/K-d_tree). Pourtant, je savais dès la lecture qu'il fallait un [*BSP*](https://en.wikipedia.org/wiki/Binary_space_partitioning) et des intersection [*aabb*](https://en.wikipedia.org/wiki/Bounding_volume). À partir de ce jour, j'étais *vraiment* dans le rouge. 

* * *
## Day 23

J'ai une première implémentation de ce problème, mais récemment, en regardant d'autres solutions, j'ai vu la représentation du board présentée ici. J'ai trouvé la forme tellement drôle que j'ai refais ma version avec!

En revanche, j'ai perdu le lien reddit, si l'auteur original se reconnaît et qu'il me le dit, je pourrais le créditer mieux! Il y a des perles en ligne...

* * *
## Day 24

J'ai craqué ce problème à la main (dans les fichiers `txt`) grâce aux encouragements d'un *teammate*. Sûrement le jour le plus éprouvant pour moi. La solution présentée ici, est une adaptation du python brillant. Je l'ai trouvé en ligne après le 25.

Je n'aurais jamais pu résoudre ce problème en machine dans les temps et je pense encore aujourd'hui qu'il n'a pas de solution générale: il s'agit de [*compréhension de programme*](https://en.wikipedia.org/wiki/Program_comprehension) et on ne sait pas bien ce que c'est.  
Par ailleurs, la [*satisfaction de contraintes*](https://en.wikipedia.org/wiki/Constraint_satisfaction_problem) du problème est triviale (c'est ce que fait la solution montrée ici).

J'ai implémenté une [*exponentiation rapide*](https://en.wikipedia.org/wiki/Exponentiation_by_squaring) pour traiter ce problème entièrement avec des `int` (et pour rigoler aussi, pas sûr du tout que ce soit plus rapide que la FPU *mais* la taille de l'entrée ne permet pas de le savoir!).

* * *
## Day 25

Ouf! Ça se finit bien: le programme utilise le `multi-buffering` et il suit linéairement le sujet du jour. La visualisation est hypnotique!

N'hésitez pas à me laisser un commentaire et happy coding!
* * *

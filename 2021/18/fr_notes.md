Le problème nous présente une variété d'[arbres binaires](https://en.wikipedia.org/wiki/Binary_tree#Internal_nodes) sur lesquelles sont définis des opérations: les *snailfish numbers*.

Ces arbres sont composés de *feuilles* qui contiennent un entier, ces feuilles sont toutes reliées entre elles (et à la *racine* de l'arbre) par des *noeuds internes* qui ne contiennent q'une unique *paires* de *liens*. Cette structure de données est aussi connue pour opérer une *classification* comme dans les [`k-d trees`](https://en.wikipedia.org/wiki/K-d_tree) ou les [`B-trees](https://en.wikipedia.org/wiki/B-tree).

Pour la lecture des entrées, j'injecte suffisamment d'espaces dans l'entrée pour pouvoir capturer *chaque* symbole séparemment en *une seule fois*.

La fonction `newPair()` a des *arguments variables*.

L'opération `explode()` tire utilement partie de la forme [*aplatie*](https://www.geeksforgeeks.org/flatten-a-binary-tree-into-linked-list/) de l'arbre pour mettre à jour des feuilles *adjacentes*; il n'y a que sous cette forme que l'information de *voisinage* est disponible.

Dans l'opération `reduce()`, j'utilise des [`drapeaux binaires`](https://en.wikipedia.org/wiki/Mask_(computing)) pour *synchroniser* le *worflow* (`done` ligne 174~184).

Enfin, j'ai profilé la version de base et découvert que l'essentiel du *runtime* de part2 consistait à attendre pour fournir des *hints* au kernel au sujet de l'utilisation *mémoire*. J'ai pris la décision de rendre la totalité de part2 concurrente: Je lance un producteur qui lance des *sous-producteurs* et en même temps, je lance des *consumers* qui calculent les magnitudes. Dans la routine `main`, je collecte et filtre les résultats. Le résultat net de cette transformation est de ramener le *runtime* de 651ms à 245ms!

MàJ Tout ce qui est dit plus haut est vrai mais trop compliqué: en repensant ce problème depuis le début, j'ai trouvé qu'en stockant les valeurs et les profondeurs des *snail numbers* la structure de donnée rendait toutes les opérations plus faciles à l'exception du calcul de *magnitude* qui a néanmoins une difficulté qu'on peut accéder. Au total j'aurais passé plus de 24h cumulés sur ce sujet. Mais c'est lui qui me fait passé la barre de la seconde sur mon mb air m1 \o/


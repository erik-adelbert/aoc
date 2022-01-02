Le problème nous présente une variété d'[arbres binaires](https://en.wikipedia.org/wiki/Binary_tree#Internal_nodes) sur lesquelles sont définis des opérations: les *snailfish numbers*.

Ces arbres sont composés de *feuilles* qui contiennent un entier, ces feuilles sont toutes reliées entre elles (et à la *racine* de l'arbre) par des *noeuds internes* qui ne contiennent q'une unique *paires* de *liens*. Cette structure de données est aussi connue pour opérer une *classification* comme dans les [`k-d trees`](https://en.wikipedia.org/wiki/K-d_tree).

L'opération `explode()` tire utilement partie de la forme [*aplatie*](https://www.geeksforgeeks.org/flatten-a-binary-tree-into-linked-list/) de l'arbre pour mettre à jour des feuilles *adjacentes*; il n'y a que sous cette forme que l'information de *voisinage* est disponible.

Dans l'opération `reduce()`, j'utilise des [`drapeaux binaires`](https://en.wikipedia.org/wiki/Mask_(computing)) pour *synchroniser* le *worflow* (`done` ligne 174~184).
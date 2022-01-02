Quand on résoud ce problème, on trouve que pour un point `p(x,y)` qui bouge pendant une [*transformation*](https://en.wikipedia.org/wiki/Transformation_matrix) autour d'un axe `a` on a:  
$$p_{n+1}.x = 2*a.x - p_{n}.x\ ||\ p_{n+1}.y = 2*a.y - p_{n}.y\\$$

On peut montrer pourquoi, à l’aide des coordonnées homogènes de *moebius*. Dans ce système, ces matrices représentent une *translation* de vecteur `u(x, y)`, une *symétrie* sur l’axe des `x` (pliure horizontale) et une sur celle des `y` (pliure verticale):  

$$\begin{pmatrix}1 & 0 & x\\0 & 1 & y\\0 & 0 & 1\\\end{pmatrix}
\begin{pmatrix}1 & 0 & 0\\0 & -1 & 0\\0 & 0 & 1\\\end{pmatrix}
\begin{pmatrix}-1 & 0 & 0\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix}$$ 

Plier la feuille en deux verticalement sur l’axe `a=x`, ça équivaut à 1) translater les points concernés de `-x` (pour les mettre dans un repère centré sur eux), 2) les symétriser verticalement et 3) les translater de `x` (pour les remettre dans le repère d’origine). Soit une combinaison des transformations suivantes: $$Tr(x).S(y).Tr(-x)$$ Ça se lit de droite à gauche, on ne peut pas changer l’ordre. Pour transformer un point, on écrit:
$$\begin{pmatrix}1 & 0 & a.x\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} . 
\begin{pmatrix}-1 & 0 & 0\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} .
\begin{pmatrix}1 & 0 & -a.x\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} .
\begin{pmatrix}p.x\\p.y\\p.z\end{pmatrix}$$

Et on voit apparaître la formule: $$p_{n+1}.x = 2*a.x - p_{n}.x$$  
Pour des problèmes où tous les points sont concernés, ça devient très puissant: on peut souvent combiner (multiplier) toutes les transformations en une seule matrice et ensuite l’appliquer en masse (`O(T+N)`) au lieu de calculer chaque transformation pour chaque point (`O(T*N)`)…
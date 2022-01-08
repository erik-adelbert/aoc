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
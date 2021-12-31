Quand on résoud ce problème, on trouve que pour un point `p` qui bouge pendant une transformation autour d'un axe `a` on a:  
$$p.x = 2*a.x - p.x\ ||\ p.y = 2*a.y - p.y\\$$

On peut montrer pourquoi, à l’aide des coordonnées homogènes de moebius:  
voici les matrices qui représentent une translation de vecteur u(x, y), une symétrie sur l’axe des x (pliure horizontale) et une sur celle des y (pliure verticale):  

$$\begin{pmatrix}1 & 0 & x\\0 & 1 & y\\0 & 0 & 1\\\end{pmatrix}
\begin{pmatrix}1 & 0 & 0\\0 & -1 & 0\\0 & 0 & 1\\\end{pmatrix}
\begin{pmatrix}-1 & 0 & 0\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix}$$ 

Plier la feuille en deux verticalement sur l’axe `x=n`, ça équivaut à 1) translater les points concernés de `-n`(pour les mettre dans un repère centré sur eux), les symétriser verticalement et les translater de `n` (pour les remettre dans le repère d’origine). Soit une combinaison des transformations suivantes: `Tr(n).Sy(x).Tr(-n)`(ça se lit de droite à gauche, on ne peut pas changer l’ordre). Pour transformer un point, on écrit:
$$\begin{pmatrix}1 & 0 & a.x\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} . 
\begin{pmatrix}-1 & 0 & 0\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} .
\begin{pmatrix}1 & 0 & -a.x\\0 & 1 & 0\\0 & 0 & 1\\\end{pmatrix} .
\begin{pmatrix}p.x\\p.y\\p.z\end{pmatrix}$$

Et on voit apparaître la formule: $$p.x = 2*a.x - p.x$$  
Pour des problèmes où tous les points sont concernés, ça devient très puissant puisqu’on combine toutes les transformations en une seule matrice et ensuite on l’applique en masse `O(N+T)` au lieu de se taper chaque transformation pour chaque point `O(N*T)`…
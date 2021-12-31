Dans ce problème, le programme doit jouer au loto. C'est à dire tenir le compte des cases qui ont été tirées sur ses cartes.
Pour modéliser une carte, je déclare une structure de donnée qui groupe:
- 5 compteurs de lignes
- 5 compteurs de colonnes
- 1 *map* des nombres et de leur position sur la carte

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

Au moment où on détecte un `bingo`, la partie est finie et il faut calculer la somme des nombres restants sur la carte: Il suffit de sommer soit les compteurs de ligne, soit ceux des colonnes.

La structure de donnée *simplifie* la résolution en supportant *facilement* et *rapidement* les opération du tirage de loto. En pratique, comme l'entrée ne contient pas de doublon, je ne vide pas la `map` dans la fonction `biff()`.
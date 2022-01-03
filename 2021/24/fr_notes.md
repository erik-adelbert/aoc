J'ai craqué ce problème à la main (dans les fichiers `txt`) grâce aux encouragements d'un *teammate*. Sûrement le jour le plus éprouvant pour moi. La solution présentée ici, est une adaptation du python brillant. Je l'ai trouvé en ligne après le 25.

Je n'aurais jamais pu résoudre ce problème en machine dans les temps et je pense encore aujourd'hui qu'il n'a pas de solution générale: il s'agit de [*compréhension de programme*](https://en.wikipedia.org/wiki/Program_comprehension) et on ne sait pas bien ce que c'est.  
Par ailleurs, la [*satisfaction de contraintes*](https://en.wikipedia.org/wiki/Constraint_satisfaction_problem) du problème est triviale (c'est ce que fait la solution montrée ici).

J'ai implémenté une [*exponentiation rapide*](https://en.wikipedia.org/wiki/Exponentiation_by_squaring) pour traiter ce problème entièrement avec des `int` (et pour rigoler aussi, pas sûr du tout que ce soit plus rapide que la FPU *mais* la taille de l'entrée ne permet pas de le savoir!).
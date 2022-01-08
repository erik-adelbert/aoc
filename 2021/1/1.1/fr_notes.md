Pour cette inauguration d'aoc2021, j'ai écrit et soumis ce programme en moins de 5mn mais ça ne m'a fait entrer que dans le top 3000, wow!

Il s'agit de compter, dans une suite de nombres, combien de fois on a une augmentation entre deux nombres successifs. Rien à dire sur le fond, ici c'est la vitesse de composition qui compte.  
En revanche, sur la forme, un petit programme comme celui-ci c'est l'occasion de sortir, du premier coup, un programme *bien écrit*. C'est à dire avec une nomenclature et un style choisis pour durer: la construction `for input.scan(){}` est un classique qu'on retrouve à l'identique dans le manuel du Go par [Donovan & Kernighan (D&K)](https://www.gopl.io). C'est effficient parce que ça énonce simplement, clairement et de manière compacte qu'elle est l'intention de programmation: on veut scanner les entrées.  
En ce qui concerne la variable `old`, elle pourrait s'appeler `last`, c'est plus traditionnel. Mais, en pratique, il m'est plus facile de manipuler des symboles d'à peu près la même longueur. Ici, on a une variable "courante" et une variable "précédente" soient: `last` (ou `prev`) et `cur` ou plus court, `old` et `cur` qui aident tout aussi bien à capturer le sens de ce qui est fait au premier coup d'oeil.  
Il y a aussi le minimalisme de `Go`, qui m'oblige à déclarer la constante `MaxInt` de façon [*idiomatique*](https://dgryski.medium.com/idiomatic-go-resources-966535376dba).



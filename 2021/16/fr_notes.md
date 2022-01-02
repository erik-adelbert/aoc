Pour résoudre ce problème linéairement, j'utilise l'excellente bibliotèque [`bitstream-go`](https://github.com/bearmini) de `bearmini` que je combine avec `math/big` de la bibliotèque standard. 

Ces bibliothèques implémentent des [interfaces](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go) standardisées comme [`Reader`](https://go.dev/tour/methods/21) et [`Writer`](https://www.grant.pizza/blog/the-beauty-of-io-writer/).

Grâce à ces interfaces, les (fonctions) filtres ont le même *prototype* (interface) que les fonctions bas niveau devant laquelle elle sont posées. Avoir des fonctions avec des prototypes identiques en entrée et en sortie permet de fabriquer des [*pipelines de données*](https://en.wikipedia.org/wiki/Pipeline_(computing)) comme sur la ligne 150.
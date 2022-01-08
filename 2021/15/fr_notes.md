Lorsqu'on utilise [*l'algorithme de Dijkstra*](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm), on ne peut plus dire grand chose d'autre: ça construit un genre de [*mycelium*](https://en.wikipedia.org/wiki/Mycelium) de données, à mesure que l'algorithme décompose (grignote) le problème. Pour accélerer le processus et guider la recherche (le grignotage), on utilise une [*priority queue*](https://en.wikipedia.org/wiki/Priority_queue). Ici, j'utilise [celle](https://pkg.go.dev/container/heap) de la bibliothèque standard de `Go`.
  
<div style="text-align:center">
  <img src="https://upload.wikimedia.org/wikipedia/commons/2/23/Dijkstras_progress_animation.gif" />
</div>

Je me méfie du code dupliqué (ou quasi-dupliqué): c'est plus long à écrire, on se trompe plus facilement en le composant, c'est pénible à maintenir. 
C'est pourquoi il y a une unique fonction `rate()`. Elle exécute une de ses deux branches en fonction de son entrée. Comme le choix est binaire, c'est un booléen qui décide son mode. C'est pour rendre ce booléen plus lisible et faciliter la mise au point que je déclare les deux constantes (o2 <- O2, co2 <- CO2). 

`rate()` mesure (sous la forme d'une `string`) les *most/least popular bits* des `inputs` (`strings` aussi) et retourne le résultat de `strconv.ParseInt()`, de la bibliothèque standard, sur cete mesure.
Au lieu de gérer ou filtrer `err` sur place (dans `rate()`), je la laisse remonter jusqu'à ce que je sois obligé de la gérer: juste avant de la transmettre sur son channel. C'est ce que je préfère comme gestion d'erreur: mes programmes ne gèrent une erreur que quand elle ne peut plus remonter. À ce moment-là, elle a souvent un sens bien défini et ça donne à une fonction filtre, le même prototype (interface) que la fonction bas niveau devant laquelle elle est posée.

Comme les problèmes sont bien séparés, je demande l'éxecution [concurrente](https://youtu.be/oV9rvDllKEg) de `rate()` dans deux `goroutines`: c'est gratuit! 



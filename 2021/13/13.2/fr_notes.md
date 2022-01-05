Ce problème déguise, sous une forme romancée, deux techniques de base, une de la [*théorie des codes*](https://en.wikipedia.org/wiki/Coding_theory) et une de [*l'informatique graphique*](https://en.wikipedia.org/wiki/Computer_graphics): Les points sont donnés sous une forme vectorielle. Après 1) les avoir transformés (*décodés*) conformément à un modèle (*code*) de pliage, il faut 2) les afficher à l'écran.

J'ai traité cet affichage comme un [*Raster Scan Display*](https://www.geeksforgeeks.org/raster-scan-displays/). Pour l'implémenter, j'utilise une [*axis aligned bounding box*](https://en.wikipedia.org/wiki/Minimum_bounding_box#Axis-aligned_minimum_bounding_box) combinée avec un [*framebuffer*](https://en.wikipedia.org/wiki/Framebuffer). Je [*rasterise*](https://en.wikipedia.org/wiki/Rasterisation) les points dans ce *buffer*. Grâce à un teammate, pendant une review, j'ai découvert que le caractère `undefined` `�` de l'`ASCII` étendu est un des plus lumineux. C'est pourquoi je le choisis comme valeur du *pixel allumé*. 

Quand j'utilise des *abstractions connues*, je gagne du temps: je facilite l'écriture, la mise au point et la maintenance du programme; je sais à l'avance ce qu'il faut faire et comment il faut le faire parce que la documentation est abondante. Si je les utilise suffisamment, je finis par les connaître par coeur.

Je trouve ce programme aussi fun qu'efficace et j'ai bien aimé l'écrire.

```
❯ cat input.txt| ./aoc13.2
�    ���  ����   �� ���    �� ���� �  �
�    �  � �       � �  �    � �    �  �
�    �  � ���     � ���     � ���  ����
�    ���  �       � �  �    � �    �  �
�    � �  �    �  � �  � �  � �    �  �
���� �  � �     ��  ���   ��  ���� �  �
```
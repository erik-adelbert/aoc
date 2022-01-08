Ce problème déguise, sous une forme romancée, deux techniques de base, une de la [*théorie des codes*](https://en.wikipedia.org/wiki/Coding_theory) et une de [*l'informatique graphique*](https://en.wikipedia.org/wiki/Computer_graphics): les points sont donnés sous une forme vectorielle. Après 1) les avoir transformés (*décodés*) conformément à un modèle (*code*) de pliage, il faut 2) les afficher à l'écran.

J'ai traité cet affichage comme un [*Raster Scan Display*](https://www.geeksforgeeks.org/raster-scan-displays/). Pour l'implémenter, j'utilise une [*aabb*](https://en.wikipedia.org/wiki/Minimum_bounding_box#Axis-aligned_minimum_bounding_box) combinée avec un [*framebuffer*](https://en.wikipedia.org/wiki/Framebuffer). Pour obtenir une *image*, je [*rasterise*](https://en.wikipedia.org/wiki/Rasterisation) les points dans ce *buffer*. Grâce à un teammate, pendant une [review](https://en.wikipedia.org/wiki/Code_review), j'ai découvert que le caractère `undefined` `�` de l'`ASCII` étendu est un des plus lumineux. C'est pourquoi je le choisis comme valeur du *pixel allumé*. 

Quand j'utilise des [*abstractions connues*](https://en.wikipedia.org/wiki/Abstraction_(computer_science)), je gagne du temps: je facilite l'écriture, la mise au point et la maintenance du programme; je sais à l'avance ce qu'il faut faire et comment il faut le faire parce que la documentation est abondante. Si je les utilise suffisamment, je finis par les connaître par coeur.

Je trouve ce problème très fun, j'ai bien aimé écrire un programme [efficace](https://en.wikipedia.org/wiki/Algorithmic_efficiency) pour le résoudre.

```
❯ cat input.txt| ./aoc13.2
�    ���  ����   �� ���    �� ���� �  �
�    �  � �       � �  �    � �    �  �
�    �  � ���     � ���     � ���  ����
�    ���  �       � �  �    � �    �  �
�    � �  �    �  � �  � �  � �    �  �
���� �  � �     ��  ���   ��  ���� �  �
```
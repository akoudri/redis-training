# RDB

save 900 1      // Sauvegarde après 900 secondes si au moins 1 clé a été modifiée
save 300 10     // Sauvegarde après 300 secondes si au moins 10 clés ont été modifiées
save 60 10000   // Sauvegarde après 60 secondes si au moins 10000 clés ont été modifiées

sudo systemctl restart redis // À exécuter à chaque modification

# AOF

appendonly yes
appendfsync everysec 
// always : Synchronise le fichier AOF à chaque écriture. C'est le plus sûr, mais aussi le plus lent.
// everysec : Synchronise le fichier AOF environ toutes les secondes. C'est un bon compromis entre sécurité et performance.
// no : Laisse le système d'exploitation décider quand synchroniser le fichier AOF, ce qui est le moins sûr mais le plus rapide.



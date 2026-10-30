export type Language = 'en'|'fr';
const fr:Record<string,string>={
'The thread between events':'Le fil entre les événements','Synthetic workspace':'Espace synthétique','Owner access':'Accès propriétaire','Passphrase':'Phrase secrète','Open workspace':'Ouvrir l’espace','Signing in…':'Connexion…','Sign out':'Déconnexion','Journey explorer':'Explorateur de parcours','Follow what happened.':'Suivez ce qui s’est passé.','Understand the path.':'Comprenez le parcours.','Observed journeys, healthy loops, and the evidence in between.':'Parcours observés, boucles normales et preuves qui les relient.','Explore':'Explorer','Configure':'Configurer','Reference views':'Vues de référence','Data lifecycle':'Cycle des données','Historical':'Historique','Live':'En direct','Version':'Version','All versions':'Toutes les versions','Definition':'Définition','All definitions':'Toutes les définitions','Anonymous visit':'Visite anonyme','All app traffic':'Tout le trafic de l’application','From':'Du','To':'Au','Apply filters':'Appliquer les filtres','Refresh':'Actualiser','Pause updates':'Suspendre les mises à jour','Resume updates':'Reprendre les mises à jour','Show seven-day app baseline':'Afficher la référence des sept jours','Clean live window':'Fenêtre directe seule','App-wide baseline':'Référence de l’application','Live observation':'Observation en direct','Paused · collection continues':'En pause · la collecte continue','Connecting':'Connexion','Connected':'Connecté','Reconnecting':'Reconnexion','Historical snapshot':'Instantané historique','Journeys':'Parcours','Completed goals':'Objectifs atteints','Active':'Actifs','Inactive · outcome unknown':'Inactifs · résultat inconnu','Unassigned records':'Enregistrements non attribués','Journey map':'Carte des parcours','Node area: median time at stage. Edge width: transition occurrences.':'Surface des nœuds : durée médiane. Épaisseur des liens : transitions observées.','Unavailable':'Indisponible','Map and timing table':'Tableau de la carte et des durées','Stage':'Étape','Occurrences':'Occurrences','Visits':'Visites','Potential stops':'Arrêts possibles','Pending':'En attente','Timing basis':'Base temporelle','Median':'Médiane','Slower cases (p90)':'Cas plus lents (p90)','Samples / exclusions':'Échantillons / exclusions','Observed transitions':'Transitions observées','Transition':'Transition','Ordered journeys':'Parcours ordonnés','No journeys in this view':'Aucun parcours dans cette vue','Receive synthetic telemetry, then check the filters and mappings.':'Recevez de la télémétrie synthétique, puis vérifiez les filtres et les correspondances.','No events received yet':'Aucun événement reçu','No step selected':'Aucune étape sélectionnée','Show all paths':'Afficher tous les chemins','Inspect paths containing':'Examiner les chemins contenant','No inferred cause':'Aucune cause déduite','Completed':'Terminé','Inactive':'Inactif','Exited':'Sortie explicite','Recovered':'Rétabli','Partial beginning':'Début partiel','Mixed versions':'Versions mixtes','Ordered evidence':'Preuves ordonnées','Event':'Événement','Occurrence time':'Heure de l’événement','Receipt time':'Heure de réception','Producer / signal':'Producteur / signal','Context':'Contexte','Copy sanitized context':'Copier le contexte filtré','Copied':'Copié','Discovered events':'Événements découverts','Event mappings':'Correspondances d’événements','Operational goal':'Objectif opérationnel','Add mapping':'Ajouter une correspondance','Preview matches':'Prévisualiser les correspondances','Activate new definition':'Activer une nouvelle définition','Source':'Source','Any source':'Toute source','Browser':'Navigateur','Server':'Serveur','Kind':'Type','Condition key':'Clé de condition','Equals':'Égale','Remove':'Supprimer','Next-event estimated timing':'Durée estimée jusqu’au prochain événement','Save a reference view':'Enregistrer une vue de référence','Name':'Nom','Mode':'Mode','Goal-backward':'Vers l’objectif','Beginning / end segment':'Segment début / fin','Origin-forward':'Depuis l’origine','Beginning event':'Événement de début','End / goal event':'Événement de fin / objectif','Save view':'Enregistrer la vue','Open view':'Ouvrir la vue','No saved views':'Aucune vue enregistrée','Captured boundaries are preserved.':'Les limites enregistrées sont préservées.','Retained summaries':'Synthèses conservées','Collection health':'État de la collecte','Limits':'Limites','Delete this visit':'Supprimer cette visite','Delete app data and disable intake':'Supprimer les données et désactiver la réception','Visit scope':'Portée de la visite','Delete visit':'Supprimer la visite','Confirmation':'Confirmation','Cancel':'Annuler','First occurrence':'Première occurrence','Last occurrence':'Dernière occurrence','Reference endpoints':'Extrémités de référence','Clear reference':'Effacer la référence','Showing':'Affichage','of':'sur','Load more evidence':'Charger plus de preuves','Timing observations':'Observations temporelles','No completed timing samples':'Aucun échantillon temporel complet','Selected app':'Application sélectionnée','Sentrame demonstration':'Démonstration Sentrame','Update available':'Mise à jour disponible','Session expired. Sign in again.':'Session expirée. Reconnectez-vous.','Request failed. Try again.':'Échec de la requête. Réessayez.','Date range must end after its beginning.':'La fin doit être après le début.','Safe diagnostic projection':'Projection diagnostique filtrée','Evidence details':'Détails des preuves','No reference endpoints in this boundary':'Aucune extrémité dans cette limite','Boundary / unavailable endpoint':'Limite / extrémité indisponible','Low support':'Faible échantillon','Sufficient support':'Échantillon suffisant','Definition activated':'Définition activée','Reference saved':'Référence enregistrée','Visit deleted':'Visite supprimée','Application disabled':'Application désactivée','Preview':'Prévisualisation','No personal identity joins':'Aucun lien d’identité personnelle','Elapsed is not attention.':'La durée écoulée ne mesure pas l’attention.','No individual drilldown after detail expiry.':'Aucun détail individuel après expiration.','Delete':'Supprimer','Live records':'Enregistrements en direct','Baseline journeys':'Parcours de référence','Local qualification only':'Qualification locale uniquement'
};
Object.assign(fr,{"Activation creates a new prospective definition. Captured journeys keep their original rule version. Browser claims cannot establish a server commitment.": "L’activation crée une nouvelle définition pour les événements futurs. Les parcours enregistrés gardent leur version d’origine. Un événement du navigateur ne peut pas prouver une validation serveur.", "Deletion and correction limits": "Limites de suppression et de correction", "Detailed evidence expires by occurrence time; supported first-receipt fallback is labeled. Reads exclude expired evidence immediately. Known replays do not extend its lifetime.": "Les preuves détaillées expirent selon l’heure de l’événement ; tout recours à la première réception est indiqué. Les preuves expirées sont exclues immédiatement. Une retransmission reconnue ne prolonge pas leur durée de vie.", "Diagnostic records, qualified occurrences and affected visits are distinct. Shared stacks group similar failures; proximity does not establish cause. Free-text messages and stacks are omitted by the default privacy projection.": "Les enregistrements diagnostiques, les occurrences établies et les visites concernées sont distincts. Une pile commune regroupe des erreurs semblables ; leur proximité ne prouve pas une cause. Les messages et piles libres sont omis par la politique de confidentialité par défaut.", "Eligible / comparable": "Admissibles / comparables", "Evidence before inference": "Les preuves avant les déductions", "Evidence, with its limits.": "Les preuves, avec leurs limites.", "Identifiable visit deletion removes retained detail and supported contributions. After detail expiry, individual summary subtraction is unavailable; explicit app/time-bucket deletion is required. Inferred journey boundaries cannot guarantee journey-only replay protection. Client copies, filesystem snapshots and backups are outside server erasure guarantees.": "La suppression d’une visite identifiable retire ses détails conservés et ses contributions prises en charge. Après expiration, la soustraction individuelle des synthèses est impossible ; une suppression explicite par application et période est nécessaire. Des limites de parcours déduites ne garantissent pas la protection contre les retransmissions pour un seul parcours. Les copies clients, instantanés et sauvegardes sont hors des garanties d’effacement du serveur.", "Layers have separate counting units; counts are not added.": "Les couches ont des unités distinctes ; leurs nombres ne sont pas additionnés.", "Matching records": "Enregistrements correspondants", "Meaning from received events": "Donner du sens aux événements reçus", "Observed path": "Chemin observé", "Observed path frequencies · rare paths remain visible": "Fréquences des chemins observés · les chemins rares restent visibles", "One row per recorded full path, definition and initial app version. Percentages use comparable non-partial started journeys; partial observations are separate. Low support is fewer than 20 eligible journeys, never a bug label.": "Une ligne par chemin complet enregistré, définition et version initiale de l’application. Les pourcentages utilisent les parcours commencés comparables et non partiels ; les observations partielles restent séparées. Moins de 20 parcours admissibles signifie un faible échantillon, jamais un défaut avéré.", "Ordered detail": "Détails ordonnés", "Potential stops count eligible journeys ending here; pending journeys are censored. p90 requires at least 20 completed samples of the same timing basis. Missing ends are never zero.": "Les arrêts possibles comptent les parcours admissibles finissant ici ; les parcours en attente sont censurés. Le p90 exige au moins 20 mesures terminées de même base temporelle. Une fin manquante ne vaut jamais zéro.", "Private by minimization": "Confidentialité par minimisation", "Quality": "Qualité", "Receipt lag": "Délai de réception", "Record / action": "Enregistrement / action", "Records in selection": "Enregistrements sélectionnés", "Return to a useful question": "Retrouver une question utile", "Safe scope available": "Contexte sûr disponible", "Share": "Part", "Standard telemetry is accepted without journey fields. Missing scope/tab, unmapped meaning or metrics receipt does not become a fabricated journey.": "La télémétrie standard est acceptée sans champs de parcours. L’absence de contexte ou d’onglet, une signification non définie et les métriques reçues ne créent pas de parcours fictif.", "Summary support is limited to declared app/day/version/definition/timing dimensions. Summary-only periods do not reconstruct paths. Approximate histogram quantiles differ from exact retained-detail samples.": "Les synthèses se limitent aux dimensions déclarées : application, jour, version, définition et durée. Une période couverte uniquement par des synthèses ne permet pas de reconstruire les chemins. Les quantiles approximatifs des histogrammes diffèrent des échantillons détaillés exacts.", "Synthetic local data only. Standard reception does not require Sentrame fields. Correlation is not identity, authority or causality. Metrics are receipt-only. Home Server, shared Tunnel and real application integration remain unqualified.": "Données synthétiques locales uniquement. La réception standard n’exige aucun champ Sentrame. La corrélation ne prouve ni identité, ni autorité, ni causalité. Les métriques ne servent qu’à inspecter la réception. Le Home Server, le Tunnel partagé et l’intégration réelle restent non qualifiés.", "Version / definition": "Version / définition", "days": "jours", "months": "mois", "step": "étape", "explicit exit": "sortie explicite", "stage start": "début d’étape", "stage end": "fin d’étape", "activity": "activité", "Measured elapsed": "Durée écoulée mesurée", "Observed active": "Activité observée", "Next-event estimated": "Estimation jusqu’au prochain événement", "Latest available timing policy": "Dernière politique temporelle disponible", "Same basis and definition across all nodes. Area is capped at 30 minutes; exact values stay visible.": "Même base et même définition pour tous les nœuds. La surface est plafonnée à 30 minutes ; les valeurs exactes restent visibles.", "Captured journey boundaries are preserved; cohort selection does not resegment. Counts describe observed eligible journeys; timing kinds and censoring remain separate.": "Les limites des parcours enregistrés sont préservées ; la sélection d’une cohorte ne les redéfinit pas. Les nombres décrivent les parcours admissibles observés ; les bases temporelles et la censure restent séparées.", "Low support: fewer than 20 eligible journeys. Rarity or a large timing node is not a defect.": "Faible échantillon : moins de 20 parcours admissibles. La rareté ou un grand nœud temporel ne prouve pas un défaut.", "Raw detail can be incomplete due to instrumentation, sampling, retention or filters; conversion describes only observed eligible journeys.": "Les détails peuvent être incomplets à cause de l’instrumentation, de l’échantillonnage, de la conservation ou des filtres ; la conversion décrit uniquement les parcours admissibles observés.", "Passphrase not accepted": "Phrase secrète refusée", "Owner login required": "Connexion propriétaire requise", "Wait one minute before retrying": "Attendez une minute avant de réessayer", "Goal must be a mapped event": "L’objectif doit correspondre à un événement défini", "Conflicting event mappings; use one deterministic mapping per event": "Correspondances contradictoires ; utilisez une seule correspondance déterministe par événement", "Choose a mapped goal and 1–100 mappings": "Choisissez un objectif défini et de 1 à 100 correspondances", "30-day detail · 12-month supported summaries": "Détails sur 30 jours · synthèses prises en charge sur 12 mois"});
Object.assign(fr,{"Evidence window": "Fenêtre des preuves", "Window clipping does not create direct transitions.": "Le découpage de la fenêtre ne crée pas de transitions directes.", "Right-censored at the selected window. Later evidence is outside this view; the captured boundary is preserved.": "Parcours censuré à la fin de la fenêtre sélectionnée. Les preuves ultérieures sont hors de cette vue ; la limite enregistrée est préservée.", "Captured status": "État enregistré", "Journey summaries finalize at captured start plus 30 days. Later corrections cannot revise a finalized summary from expired individual detail.": "Les synthèses de parcours deviennent définitives 30 jours après leur début enregistré. Les corrections ultérieures ne peuvent pas réviser une synthèse définitive à partir de détails individuels expirés.", "Finalized summaries are app-wide and use whole UTC days. Individual visit scope is unsupported after detail expiry; select all app traffic and whole-day bounds to query summaries.": "Les synthèses définitives couvrent toute l’application et des jours UTC complets. Une visite individuelle n’est plus consultable après expiration des détails ; sélectionnez tout le trafic et des journées entières pour consulter les synthèses."});
Object.assign(fr,{"Completed samples / exclusions": "Mesures complètes / exclusions", "Delete one summary UTC day": "Supprimer un jour UTC de synthèse", "This removes only the selected app-wide summary day. Retained detail, other days and intake are unchanged. It does not identify an individual contribution.": "Cette action supprime uniquement le jour de synthèse sélectionné pour toute l’application. Les détails conservés, les autres jours et la réception restent inchangés. Elle n’identifie aucune contribution individuelle.", "Summary UTC day": "Jour UTC de synthèse", "Delete summary day": "Supprimer le jour de synthèse", "Summary day deleted": "Jour de synthèse supprimé"});
Object.assign(fr,{"Records in app retention": "Enregistrements conservés de l’application", "Discovery and available choices use bounded app-wide retained evidence, independently of the current live or visit filters.": "La découverte et les choix disponibles utilisent les preuves conservées de toute l’application, dans une limite définie et indépendamment des filtres de visite ou de direct."});
Object.assign(fr,{"Beginning occurrence": "Occurrence de début", "End / goal occurrence": "Occurrence de fin / objectif"});
Object.assign(fr,{
  "Journey workspace": "Espace des parcours",
  "Graph view": "Vue du graphe",
  "Hierarchy": "Hiérarchie",
  "Force": "Forces",
  "Sankey": "Sankey",
  "Graph highlight": "Mise en évidence du graphe",
  "Highlights": "Mises en évidence",
  "Rare paths": "Chemins rares",
  "Endings": "Fins de parcours",
  "Warnings & exceptions": "Avertissements et exceptions",
  "Recovery": "Rétablissement",
  "References": "Références",
  "Workspace menu": "Menu de l’espace",
  "Hide discovered events": "Masquer les événements découverts",
  "Show discovered events": "Afficher les événements découverts",
  "Unmapped events": "Événements sans correspondance",
  "Mappings & goal": "Correspondances et objectif",
  "Save current view…": "Enregistrer la vue actuelle…",
  "About this evidence": "À propos de ces preuves",
  "Select an end stage": "Sélectionnez une étape de fin",
  "Seven-day baseline": "Référence des sept jours",
  "Selection details": "Détails de la sélection",
  "Close inspector": "Fermer l’inspecteur",
  "Observed transition": "Transition observée",
  "Discovered event": "Événement découvert",
  "Journey stage": "Étape du parcours",
  "Current selection": "Sélection actuelle",
  "Baseline only": "Référence uniquement",
  "A transition is observed adjacency, not proof of cause. Open journeys to inspect its order and context.": "Une transition indique une succession observée, sans prouver une cause. Ouvrez les parcours pour examiner leur ordre et leur contexte.",
  "Show journeys": "Afficher les parcours",
  "Inspect beginning": "Examiner le début",
  "Inspect end": "Examiner la fin",
  "Unmapped event": "Événement sans correspondance",
  "Received telemetry without an assigned journey meaning. No relationship is inferred from proximity.": "Télémétrie reçue sans signification définie pour les parcours. La proximité ne permet pas de déduire une relation.",
  "Configure event…": "Configurer l’événement…",
  "Configure event": "Configurer l’événement",
  "Silence does not establish failure.": "Le silence ne prouve pas un échec.",
  "Incoming": "Entrées",
  "Outgoing": "Sorties",
  "Repeated here": "Répétition ici",
  "No observed transitions in this selection": "Aucune transition observée dans cette sélection",
  "transitions": "transitions",
  "Time at stage": "Durée à cette étape",
  "samples": "échantillons",
  "excluded": "exclus",
  "complete": "complets",
  "p90 requires 20 completed samples of the same basis. Missing ends are not zero. Elapsed time is not attention.": "Le p90 exige 20 mesures terminées de même base temporelle. Une fin manquante ne vaut pas zéro. Le temps écoulé ne mesure pas l’attention.",
  "Events in this stage": "Événements de cette étape",
  "Paths into here…": "Chemins menant ici…",
  "Paths from here…": "Chemins partant d’ici…",
  "Use as beginning…": "Utiliser comme début…",
  "This selection has no retained evidence in the current filters.": "Aucune preuve conservée pour cette sélection avec les filtres actuels.",
  "Table": "Tableau",
  "Journeys & evidence": "Parcours et preuves",
  "Rare: at most 10% of comparable complete paths, with at least 20 closed journeys per definition and version. Low support is not a defect.": "Rareté : au plus 10 % des chemins complets comparables, avec au moins 20 parcours clos par définition et version. Un faible échantillon ne prouve pas un défaut.",
  "Path counts": "Nombres par chemin",
  "Endings distinguish completed, inactive, explicit exit and pending evidence.": "Les fins de parcours distinguent les objectifs atteints, l’inactivité, les sorties explicites et les preuves en attente.",
  "Recovered journeys retain their earlier warnings and later continuation.": "Les parcours rétablis conservent leurs avertissements antérieurs et leur continuation ultérieure.",
  "Warnings and exceptions highlight linked journeys; proximity is not cause.": "Les avertissements et exceptions mettent en évidence les parcours liés ; la proximité ne prouve pas une cause.",
  "Inspect evidence": "Examiner les preuves",
  "Query work limit reached. Counts are partial; narrow your filters.": "Limite de traitement de la requête atteinte. Les nombres sont partiels ; précisez vos filtres.",
  "Set operational goal": "Définir l’objectif opérationnel",
  "Give this received event a journey meaning. Relationships come from ordered evidence, not this configuration.": "Définissez la signification de cet événement reçu dans les parcours. Les relations proviennent des preuves ordonnées, pas de cette configuration.",
  "Trusted producer": "Producteur de confiance",
  "Received examples": "Exemples reçus",
  "No retained examples in the current selection.": "Aucun exemple conservé dans la sélection actuelle.",
  "This explicitly changes the operational goal for future journeys.": "Cette action change explicitement l’objectif opérationnel des futurs parcours.",
  "without safe journey context": "sans contexte de parcours sûr",
  "Time & filters": "Période et filtres",
  "Observed paths": "Chemins observés",
  "bounds fixed": "limites fixes",
  "Traffic counts; timing remains separate by basis and definition.": "Nombres de trafic ; les durées restent séparées par base temporelle et définition.",
  "Current workspace": "Espace actuel",
  "All captured paths": "Tous les chemins enregistrés",
  "Discovery uses retained app-wide evidence. Unconnected events have no supported journey relationship.": "La découverte utilise les preuves conservées de toute l’application. Les événements non reliés n’ont aucune relation de parcours établie.",
  "Show discovered events on canvas": "Afficher les événements découverts sur le canevas",
  "records": "enregistrements",
  "No unmapped events": "Aucun événement sans correspondance",
  "Active definition": "Définition active",
  "Exploring paths to a stage does not change the operational goal.": "Explorer les chemins vers une étape ne change pas l’objectif opérationnel.",
  "Timing policy": "Politique temporelle",
  "Graph navigation": "Navigation dans le graphe",
  "Zoom out": "Réduire le zoom",
  "Zoom in": "Agrandir le zoom",
  "Fit": "Ajuster",
  "Receive events or adjust the time and visit filters.": "Recevez des événements ou ajustez les filtres de période et de visite.",
  "Bar height and ribbon width: occurrences": "Hauteur des barres et largeur des rubans : occurrences",
  "Arrows: observed direction · line width: traffic": "Flèches : direction observée · épaisseur des liens : trafic",
  "Muted: seven-day baseline · Green: live · counts stay separate": "Atténué : référence des sept jours · Vert : direct · nombres séparés",
  "Dashed outline: unmapped event": "Contour pointillé : événement sans correspondance",
  "Rare: ≤10% of ≥20 closed, comparable journeys; partial / ongoing paths excluded. Highlighted links may also occur in common paths.": "Rareté : ≤10 % d’au moins 20 parcours clos comparables ; chemins partiels ou en cours exclus. Les liens mis en évidence peuvent aussi appartenir à des chemins courants.",
  "Tab to select stages and connections. Enter opens the inspector. On the canvas use arrow keys to pan, plus and minus to zoom, zero to fit. Drag empty canvas to pan. Timing and ordered evidence are in the inspector and table.": "Utilisez Tab pour sélectionner les étapes et les liens, puis Entrée pour ouvrir l’inspecteur. Sur le canevas, les flèches déplacent la vue, plus et moins règlent le zoom, et zéro ajuste la vue. Faites glisser une zone vide pour déplacer la vue. Les durées et les preuves ordonnées sont dans l’inspecteur et le tableau.",
  "occurrences": "occurrences",
  "baseline": "référence",
  "live": "direct",
  "unmapped event": "événement sans correspondance",
  "operational goal": "objectif opérationnel",
  "Unmapped": "Sans correspondance",
  "Unmapped event — no journey relationship inferred": "Événement sans correspondance — aucune relation de parcours déduite",
  "return / repeat path": "chemin de retour / répétition",
  "End": "Fin",
  "Rare": "Rare",
  "Showing nodes": "Nœuds affichés",
  "Narrow filters; omitted relationships are not joined.": "Précisez les filtres ; les relations omises ne sont pas reliées.",
  "Showing": "Affichage",
  "nodes. Narrow filters; omitted relationships are not joined.": "nœuds. Précisez les filtres ; les relations omises ne sont pas reliées.",
  "Close": "Fermer",
  "Overview": "Vue d’ensemble",
  "Workspace overview": "Vue d’ensemble de l’espace",
  "Whole workspace": "Espace complet",
  "Save reference…": "Enregistrer une référence…",
  "Preview paths": "Prévisualiser les chemins",
  "No matching journeys": "Aucun parcours correspondant",
  "Save reference": "Enregistrer la référence"
});
Object.assign(fr,{'Mixed endings':'Fins de parcours mixtes','baseline transitions':'transitions de référence','Unknown':'Inconnu'});
Object.assign(fr,{
  "One row per observed sequence, definition and initial app version. Rarity uses closed, non-partial, uncensored journeys; at least 20 comparable journeys and at most 10% sequence share. Partial and ongoing observations remain in total counts, not the rarity denominator.": "Une ligne par séquence observée, définition et version initiale de l’application. La rareté utilise les parcours clos, non partiels et non censurés ; au moins 20 parcours comparables et une part de séquence d’au plus 10 %. Les observations partielles ou en cours restent dans les nombres totaux, sans entrer dans le dénominateur de rareté.",
  "Choose first or last endpoint occurrence": "Choisissez la première ou la dernière occurrence de l’extrémité",
  "Choose a supported reference mode": "Choisissez un mode de référence pris en charge",
  "Select reference endpoints": "Sélectionnez les extrémités de référence",
  "Choose a saved reference or an unsaved preview, not both": "Choisissez une référence enregistrée ou une prévisualisation non enregistrée, pas les deux",
  "Name required": "Nom requis",
  "Choose a supported graph layout": "Choisissez une disposition de graphe prise en charge",
  "Choose a supported graph highlight": "Choisissez une mise en évidence de graphe prise en charge"
});
Object.assign(fr,{'Clear selected visit':'Effacer la visite sélectionnée'});
Object.assign(fr,{
 'Paths into here':'Chemins vers cette étape','Paths from here':'Chemins depuis cette étape','Between stages':'Entre deux étapes',
 'Saves this live observation window with fixed bounds. The seven-day baseline is a separate comparison layer.':'Enregistre cette fenêtre d’observation en direct avec des limites fixes. La référence de sept jours reste une couche de comparaison distincte.'
});
Object.assign(fr,{'Stop mapping this event':'Ne plus associer cet événement','Choose another operational goal before removing this mapping.':'Choisissez un autre objectif opérationnel avant de retirer cette correspondance.','This event will no longer be mapped in the next definition.':'Cet événement ne sera plus associé dans la prochaine définition.'});
Object.assign(fr,{'Tab to select stages and connections. Enter opens the inspector. On the canvas use arrow keys to pan, plus and minus to zoom, zero to fit. Drag or two-finger scroll to pan; pinch to zoom. Timing and ordered evidence are in the inspector and table.':'Utilisez Tab pour sélectionner une étape ou une connexion, puis Entrée pour ouvrir les détails. Sur le canevas, les flèches déplacent la vue, plus et moins changent le zoom et zéro ajuste la vue. Glissez ou défilez à deux doigts pour déplacer la vue ; pincez pour zoomer. Les durées et les preuves ordonnées restent dans les détails et le tableau.'});
Object.assign(fr,{'Visit / operation':'Visite / opération','Error evidence on this connection':'Preuve d’erreur sur cette connexion','Error evidence attached':'Preuve d’erreur associée','Red connections: error evidence':'Connexions rouges : preuves d’erreur'});
Object.assign(fr,{
  "Visual time filter": "Filtre temporel visuel",
  "Time window": "Période sélectionnée",
  "records in highlighted bins": "enregistrements dans les intervalles sélectionnés",
  "Retained records by occurrence time": "Enregistrements conservés selon l’heure de l’événement",
  "records in context": "enregistrements dans la période affichée",
  "Maximum bin": "Maximum par intervalle",
  "Start": "Début",
  "Drag the bars to select a range, or use the Start and End sliders. Applies when released.": "Faites glisser les barres pour sélectionner une période, ou utilisez les curseurs Début et Fin. La sélection s’applique au relâchement.",
  "Counts are retained records, not journeys or people.": "Les nombres concernent les enregistrements conservés, pas les parcours ni les personnes.",
  "The selected dates are outside this context. The graph keeps its existing range until you choose a new one.": "Les dates sélectionnées sont hors de la période affichée. Le graphe conserve sa période actuelle jusqu’à votre prochaine sélection.",
  "The existing range extends beyond this context. The handles show its visible portion; choosing a new range replaces it.": "La période actuelle dépasse celle du graphique. Les curseurs montrent sa partie visible ; une nouvelle sélection la remplace.",
  "Partial timeline: the retained-record query reached its work limit. Narrow the context or filters.": "Histogramme partiel : la requête sur les enregistrements conservés a atteint sa limite de traitement. Réduisez la période ou précisez les filtres.",
  "Timeline data table": "Tableau des données temporelles",
  "No retained records timeline is available in this context.": "Aucun histogramme des enregistrements conservés n’est disponible pour cette période.",
  "Loading time context…": "Chargement de la période…",
  "Records": "Enregistrements",
  "Partial": "Partiel",
  "Advanced: exact times or correlation ID": "Options avancées : heures précises ou identifiant de corrélation",
  "Chart context": "Période du graphique",
  "Correlation ID (optional)": "Identifiant de corrélation (facultatif)",
  "Correlation details": "Détails de corrélation",
  "Last 15 minutes": "15 dernières minutes",
  "Last 24 hours": "24 dernières heures",
  "Last hour": "Dernière heure",
  "Loading…": "Chargement…",
  "Paste a known session or operation ID": "Collez un identifiant connu de session ou d’opération",
  "Quick range": "Périodes rapides",
  "This code groups related events; it does not identify a person or link separate visits.": "Ce code regroupe des événements liés ; il n’identifie pas une personne et ne relie pas des visites distinctes.",
  "Timeline context": "Période de l’histogramme",
  "To follow one visit, select a node, open its journeys and choose Follow this visit. No session ID is needed.": "Pour suivre une visite, sélectionnez un nœud, ouvrez ses parcours et choisissez Suivre cette visite. Aucun identifiant de session n’est nécessaire.",
  "Follow this visit": "Suivre cette visite",
  "Follow this operation": "Suivre cette opération",
  "Visit / operation": "Visite / opération",
  "Clear selected visit": "Effacer la visite sélectionnée",
  "Paths into here": "Chemins menant ici",
  "Paths from here": "Chemins partant d’ici",
  "Between stages": "Entre deux étapes",
  "Timeline unavailable": "Histogramme indisponible",
  "Choose timeline window 1h, 24h, 7d or 30d": "Choisissez une période de 1 h, 24 h, 7 j ou 30 j",
  "Choose a positive configuration definition version": "Choisissez un numéro positif de version de configuration",
  "Timeline bin unavailable": "Intervalle de l’histogramme indisponible",
  "1h": "1 h",
  "24h": "24 h",
  "7d": "7 j",
  "30d": "30 j"
});
Object.assign(fr,{'Error evidence':'Preuve d’erreur'});
Object.assign(fr,{
  "Unable to load errors.": "Impossible de charger les erreurs.",
  "Copy unavailable. Select the displayed sanitized text.": "Copie indisponible. Sélectionnez le texte filtré affiché.",
  "Errors & exceptions": "Erreurs et exceptions",
  "Discover errors from the application, including background work and events without a visit. You do not need to be the person using the app.": "Découvrez les erreurs de l’application, y compris celles des tâches en arrière-plan et des événements sans visite. Vous n’avez pas besoin d’être la personne qui utilise l’application.",
  "Whole application": "Toute l’application",
  "Current graph filters": "Filtres actuels de la carte",
  "Find an error": "Rechercher une erreur",
  "Message, exception, service or trace ID": "Message, exception, service ou identifiant de trace",
  "Search": "Rechercher",
  "diagnostic records": "enregistrements diagnostiques",
  "Selected graph context": "Contexte sélectionné sur la carte",
  "Retained detail · last 30 days": "Détails conservés · 30 derniers jours",
  "Results are partial. Narrow the search or graph filters to inspect more evidence.": "Les résultats sont partiels. Précisez la recherche ou les filtres de la carte pour examiner davantage de preuves.",
  "Unknown service": "Service inconnu",
  "Unknown version": "Version inconnue",
  "Message": "Message",
  "Stack trace": "Pile d’appels",
  "No retained stack trace was supplied for this record.": "Aucune pile d’appels conservée n’a été fournie pour cet enregistrement.",
  "Diagnostic text was shortened to the retention limits.": "Le texte diagnostique a été raccourci selon les limites de conservation.",
  "Producer": "Producteur",
  "Explore related context": "Explorer le contexte associé",
  "Retained diagnostic fields": "Champs diagnostiques conservés",
  "No matching errors": "Aucune erreur correspondante",
  "Try the whole application or a different search. No records is not proof that no errors occurred.": "Essayez toute l’application ou une autre recherche. L’absence d’enregistrements ne prouve pas qu’aucune erreur ne s’est produite.",
  "Messages and stacks are bounded and scrubbed before storage. Scrubbing is imperfect; emit controlled diagnostic text. Detail expires after 30 days. Record counts are not distinct incidents or people.": "Les messages et les piles sont limités et filtrés avant leur stockage. Le filtrage est imparfait : émettez un texte diagnostique maîtrisé. Les détails expirent après 30 jours. Le nombre d’enregistrements ne correspond pas à un nombre d’incidents distincts ou de personnes.",
  "Diagnostic retention policy": "Politique de conservation des diagnostics",
  "Fatal": "Fatale",
  "Error": "Erreur",
  "Exception": "Exception",
  "Correlation context is available; it does not establish a person or a cause.": "Un contexte de corrélation est disponible ; il ne permet pas d’identifier une personne ni d’établir une cause.",
  "No visit correlation was supplied. This error remains discoverable without inventing a journey.": "Aucune corrélation de visite n’a été fournie. Cette erreur reste accessible sans inventer un parcours.",
  "Errors": "Erreurs",
  "Error evidence on this connection": "Preuve d’erreur sur ce lien",
  "Error evidence attached": "Preuve d’erreur associée",
  "Red connections: error evidence": "Liens rouges : preuves d’erreur",
  "Error evidence on this stage": "Preuve d’erreur sur cette étape"
});
Object.assign(fr,{'No stack trace was retained for this record.':'Aucune pile d’appels n’a été conservée pour cet enregistrement.','Sanitized stack as emitted. Source maps have not been applied.':'Pile filtrée telle qu’émise. Les cartes de sources n’ont pas été appliquées.','Earlier records may lack details omitted by the previous privacy policy.':'Les anciens enregistrements peuvent ne pas contenir les détails omis par l’ancienne politique de confidentialité.','Shows error records inside the selected dates; a journey can include later follow-up.':'Affiche les erreurs dans la période sélectionnée ; un parcours peut inclure un suivi ultérieur.'});
Object.assign(fr,{'Error record':'Enregistrement d’erreur'});
export function translate(key:string,lang:Language):string{return lang==='fr'?(fr[key]??key):key;}

Object.assign(fr,{
  'Reset layout':'Réinitialiser',
  'Tab to select stages and connections. Enter opens the inspector. On the canvas use arrow keys to pan, plus and minus to zoom, zero to fit. Drag nodes to rearrange; Shift plus arrow keys also moves nodes. Hierarchy and Sankey keep columns fixed. Drag empty space or two-finger scroll to pan; pinch to zoom. Timing and ordered evidence are in the inspector and table.':'Utilisez Tab pour sélectionner les étapes et les liens, puis Entrée pour ouvrir les détails. Sur le canevas, les flèches déplacent la vue, plus et moins changent le zoom et zéro ajuste la vue. Faites glisser les nœuds ou utilisez Maj et les flèches pour les déplacer. Hiérarchie et Sankey gardent leurs colonnes fixes. Glissez une zone vide ou défilez à deux doigts pour déplacer la vue ; pincez pour zoomer. Les durées et les preuves ordonnées restent dans les détails et le tableau.'
});

Object.assign(fr,{
  "Help: Stage": "Aide : Étape",
  "Help: Kind": "Aide : Type",
  "Help: Trusted producer": "Aide : Producteur de confiance",
  "Help: Condition key": "Aide : Clé de condition",
  "Help: Equals": "Aide : Égal à",
  "The name shown on the journey graph. Use a short, meaningful label such as Checkout or Vote saved. Several events can share a stage name.": "Le nom affiché sur le graphe du parcours. Choisissez un libellé court et clair, comme Paiement ou Vote enregistré. Plusieurs événements peuvent partager le même nom d’étape.",
  "For stage start and stage end, use the same stage name. For activity, this name does not add a graph node.": "Pour le début et la fin d’étape, utilisez le même nom. Pour une activité, ce nom n’ajoute pas de nœud au graphe.",
  "Step: adds an occurrence to the journey path. Use this for a meaningful action or milestone, such as Vote saved.": "Étape : ajoute une occurrence au parcours. Utilisez ce type pour une action ou un jalon significatif, comme Vote enregistré.",
  "Explicit exit: ends the journey with an emitted reason, such as recipes exhausted. It does not automatically mean failure.": "Sortie explicite : termine le parcours avec une raison émise, comme recettes épuisées. Cela ne signifie pas automatiquement un échec.",
  "Stage start / stage end: mark entry and exit for elapsed timing. Pair the same stage name and emitted occurrence ID from the same producer and service. The start adds a graph step; the end does not.": "Début / fin d’étape : marquent l’entrée et la sortie pour mesurer la durée écoulée. Associez le même nom d’étape et le même identifiant d’occurrence émis par le même producteur et service. Le début ajoute une étape au graphe ; la fin n’en ajoute pas.",
  "Activity: keeps supporting evidence without adding a graph step. With activity.state set to active, visible, idle or hidden, it can support observed active-time calculations; it does not measure attention.": "Activité : conserve des preuves complémentaires sans ajouter d’étape au graphe. Avec activity.state défini sur active, visible, idle ou hidden, elle peut contribuer au calcul du temps actif observé ; elle ne mesure pas l’attention.",
  "All mapped kinds refresh the inactivity timeout. The operational goal is selected separately.": "Tous les types configurés repoussent le délai d’inactivité. L’objectif opérationnel est choisi séparément.",
  "Choose which configured producer may match this event: Browser, Server, or Any source. Use Server for an event that proves a saved change.": "Choisissez le producteur configuré autorisé pour cet événement : Navigateur, Serveur ou Toute source. Utilisez Serveur pour un événement prouvant un changement enregistré.",
  "This uses the receiver’s trusted producer assignment. Selecting Server does not turn a browser claim into trusted server evidence.": "Ce choix utilise l’attribution de confiance du récepteur. Sélectionner Serveur ne transforme pas une déclaration du navigateur en preuve serveur de confiance.",
  "Optional: match only events with a particular retained attribute. Leave blank to match by event name and producer alone.": "Facultatif : ne retenir que les événements ayant un attribut conservé particulier. Laissez vide pour utiliser uniquement le nom de l’événement et le producteur.",
  "Use a safe attribute key from Received examples, such as reason. Pair it with Equals. Attributes removed by privacy filtering cannot be used.": "Utilisez une clé d’attribut autorisée des Exemples reçus, comme reason. Associez-la à Égal à. Les attributs supprimés par le filtrage de confidentialité ne peuvent pas être utilisés.",
  "The exact, case-sensitive value required for the Condition key. For example: reason equals recipes_exhausted.": "La valeur exacte attendue pour la clé de condition, en respectant les majuscules et minuscules. Par exemple : reason égal à recipes_exhausted.",
  "This is a literal match, not a search or pattern. Without a Condition key, this value has no effect. Use Preview matches before activating.": "La comparaison est littérale, sans recherche ni motif. Sans clé de condition, cette valeur n’a aucun effet. Prévisualisez les correspondances avant d’activer."
});

Object.assign(fr,{
  "This field names a piece of information inside the event. Equals is the value that information must contain. These two fields work together to filter this event’s mapping.": "Ce champ nomme une information contenue dans l’événement. Égal à indique la valeur que cette information doit avoir. Ensemble, ces deux champs filtrent cette règle de correspondance.",
  "Example: for journey.exited with {\"reason\":\"exhausted\"}, enter reason in Condition key and exhausted in Equals. Only journey.exited events with that reason match; other reasons do not.": "Exemple : pour journey.exited contenant {\"reason\":\"exhausted\"}, saisissez reason dans Clé de condition et exhausted dans Égal à. Seuls les événements journey.exited ayant cette raison correspondent ; les autres raisons sont exclues.",
  "Another example: for completion.entered with {\"reason\":\"agreed\"}, enter reason in Condition key and agreed in Equals.": "Autre exemple : pour completion.entered contenant {\"reason\":\"agreed\"}, saisissez reason dans Clé de condition et agreed dans Égal à.",
  "Open Received examples below and look inside the event attributes. Copy the attribute name into Condition key and its value into Equals, without quotes or JSON braces. Only retained, allowed attributes can be used.": "Ouvrez Exemples reçus ci-dessous et regardez les attributs de l’événement. Copiez le nom de l’attribut dans Clé de condition et sa valeur dans Égal à, sans guillemets ni accolades JSON. Seuls les attributs conservés et autorisés peuvent être utilisés.",
  "If you want every occurrence of this event from the selected producer, leave BOTH fields empty. You do not need a condition for a normal mapping.": "Pour inclure toutes les occurrences de cet événement du producteur choisi, laissez LES DEUX champs vides. Une correspondance normale ne nécessite pas de condition.",
  "Enter the exact value you want to allow for Condition key. It is text, not a formula or search.": "Saisissez la valeur exacte attendue pour la clé de condition. Il s’agit de texte, pas d’une formule ni d’une recherche.",
  "Example: Condition key = reason; Equals = exhausted. An event containing {\"reason\":\"exhausted\"} matches. {\"reason\":\"agreed\"} does not. An event without reason does not match this example either.": "Exemple : Clé de condition = reason ; Égal à = exhausted. Un événement contenant {\"reason\":\"exhausted\"} correspond. {\"reason\":\"agreed\"} ne correspond pas. Un événement sans reason ne correspond pas non plus à cet exemple.",
  "Example: Condition key = activity.state; Equals = hidden. Only events with activity.state set to hidden match this mapping.": "Exemple : Clé de condition = activity.state ; Égal à = hidden. Seuls les événements dont activity.state vaut hidden correspondent à cette règle.",
  "Type exhausted, not reason=exhausted, not \"exhausted\", and not the whole JSON object. Capitalization matters: Exhausted is different from exhausted.": "Saisissez exhausted, et non reason=exhausted, ni \"exhausted\", ni l’objet JSON entier. Les majuscules comptent : Exhausted diffère de exhausted.",
  "Leave both fields empty if you do not need this filter. Then use Preview matches to check your rule before activating it.": "Laissez les deux champs vides si ce filtre est inutile. Prévisualisez ensuite les correspondances pour vérifier la règle avant de l’activer."
});

Object.assign(fr,{
  "Explains data completeness, privacy filtering, view freshness and verified transport support. These are reported policies and known limitations, not a live guarantee that every event arrived.": "Explique la complétude des données, le filtrage de confidentialité, la fraîcheur de la vue et les transports vérifiés. Il s’agit des politiques déclarées et des limites connues, sans garantie en temps réel que chaque événement est arrivé.",
  "Caps how many records one view can process and return to keep the app responsive. If results are truncated, counts are partial; narrow your time range or filters. These caps do not change data retention.": "Plafonne le nombre d’enregistrements qu’une vue peut traiter et renvoyer pour garder l’application réactive. Si les résultats sont tronqués, les comptages sont partiels ; réduisez la période ou les filtres. Ces plafonds ne changent pas la conservation des données."
});

Object.assign(fr,{"Highlight":"Mettre en évidence","None":"Aucune"});

Object.assign(fr,{
 'Delete reference':'Supprimer la référence',
 'Reference deleted':'Référence supprimée',
 'Only this saved view is removed. Collected events, journeys and other references are kept.':'Seule cette vue enregistrée est supprimée. Les événements collectés, les parcours et les autres références sont conservés.'
});

Object.assign(fr,{
  "Temporary viewers": "Lecteurs temporaires",
  "Temporary viewer": "Lecteur temporaire",
  "Invite someone to browse the whole application in read-only mode. They can view retained historical and live journeys, references, errors and sanitized diagnostic details. They cannot configure, save, delete or invite anyone.": "Invitez une personne à consulter toute l’application en lecture seule. Elle pourra voir les parcours historiques conservés et en direct, les références, les erreurs et les détails de diagnostic filtrés. Elle ne pourra ni configurer, ni enregistrer, ni supprimer, ni inviter une autre personne.",
  "Viewer label (optional)": "Libellé du lecteur (facultatif)",
  "For example: Product review": "Par exemple : revue du produit",
  "Access expires after": "L’accès expire après",
  "1 hour": "1 heure",
  "24 hours": "24 heures",
  "7 days": "7 jours",
  "Expiry starts when the link is created. Anyone with the link can redeem it once; send it privately. Revocation blocks new requests immediately and closes live streams within five seconds. Already copied data cannot be recalled.": "Le délai commence à la création du lien. Toute personne disposant du lien peut l’utiliser une fois ; envoyez-le en privé. La révocation bloque immédiatement les nouvelles requêtes et ferme les flux en direct sous cinq secondes. Les données déjà copiées ne peuvent pas être récupérées.",
  "Create invitation link": "Créer un lien d’invitation",
  "Copy unavailable. Select the displayed link.": "Copie indisponible. Sélectionnez le lien affiché.",
  "Invitation link": "Lien d’invitation",
  "Copy this link now": "Copiez ce lien maintenant",
  "This link is shown only here and is not saved in browser storage. The local 127.0.0.1 address works only on this Mac; this does not expose Sentrame to another computer.": "Ce lien est affiché uniquement ici et n’est pas enregistré dans le stockage du navigateur. L’adresse locale 127.0.0.1 fonctionne uniquement sur ce Mac ; cela n’expose pas Sentrame à un autre ordinateur.",
  "Copy invitation link": "Copier le lien d’invitation",
  "Invitation access": "Accès par invitation",
  "Revoked": "Révoqué",
  "Expired": "Expiré",
  "Redeemed": "Utilisé",
  "Waiting for viewer": "En attente du lecteur",
  "Expires": "Expiration",
  "Revoke access": "Révoquer l’accès",
  "No invitations yet": "Aucune invitation pour le moment",
  "Viewer sessions end on expiry, revocation or service restart. A redeemed link cannot be reused; create a new invitation if access is needed again.": "Les sessions des lecteurs se terminent à l’expiration, à la révocation ou au redémarrage du service. Un lien utilisé ne peut pas resservir ; créez une nouvelle invitation si un nouvel accès est nécessaire.",
  "Viewer access ended. Ask the owner for a new invitation.": "L’accès du lecteur est terminé. Demandez une nouvelle invitation au propriétaire.",
  "This invitation is unavailable, expired, already used or revoked. Ask the owner for a new link.": "Cette invitation est indisponible, expirée, déjà utilisée ou révoquée. Demandez un nouveau lien au propriétaire.",
  "Read-only access": "Accès en lecture seule",
  "Explore together.": "Explorez ensemble.",
  "Read-only access.": "Accès en lecture seule.",
  "This invitation grants temporary access to the whole application: historical and live journeys, references, errors and sanitized diagnostics. It does not allow changes.": "Cette invitation donne un accès temporaire à toute l’application : parcours historiques et en direct, références, erreurs et diagnostics filtrés. Elle ne permet aucune modification.",
  "Temporary invitation": "Invitation temporaire",
  "You are signed in as the owner. Sign out before accepting a viewer invitation, or open the link in a separate browser profile.": "Vous êtes connecté comme propriétaire. Déconnectez-vous avant d’accepter une invitation de lecteur, ou ouvrez le lien dans un autre profil de navigateur.",
  "The link can be redeemed once. The owner can revoke access at any time.": "Le lien peut être utilisé une seule fois. Le propriétaire peut révoquer l’accès à tout moment.",
  "Opening…": "Ouverture…",
  "Checking access…": "Vérification de l’accès…",
  "Open read-only workspace": "Ouvrir l’espace en lecture seule",
  "Viewer access ended.": "L’accès du lecteur est terminé.",
  "Your invitation expired, was revoked, or its session ended. Ask the owner for a new invitation link.": "Votre invitation a expiré, a été révoquée ou sa session est terminée. Demandez un nouveau lien d’invitation au propriétaire.",
  "Access is closed": "L’accès est fermé",
  "The workspace has been cleared from this page. Previously copied information cannot be recalled.": "L’espace de travail a été effacé de cette page. Les informations copiées auparavant ne peuvent pas être récupérées.",
  "Owner sign in": "Connexion du propriétaire"
});

Object.assign(fr,{
 'This link is shown only here and is not saved in browser storage.':'Ce lien est affiché uniquement ici et n’est pas enregistré dans le stockage du navigateur.',
 'This loopback address works only on this Mac; this does not expose Sentrame to another computer.':'Cette adresse de boucle locale fonctionne uniquement sur ce Mac ; cela n’expose pas Sentrame à un autre ordinateur.',
 'Invitation unavailable; ask the owner for a new link':'Invitation indisponible ; demandez un nouveau lien au propriétaire',
 'Read-only viewer access':'Accès lecteur en lecture seule',
 'Session expired or access unavailable. Sign in or ask the owner for a new invite.':'Session expirée ou accès indisponible. Connectez-vous ou demandez une nouvelle invitation au propriétaire.',
 'Choose 1, 24 or 168 hours':'Choisissez 1, 24 ou 168 heures',
 'Active invitation limit 100':'Limite de 100 invitations actives',
 'Viewer session limit; retry later':'Limite de sessions de lecteurs atteinte ; réessayez plus tard',
 'Invitation busy; retry':'Invitations occupées ; réessayez',
 'Application is disabled':'L’application est désactivée'
});

Object.assign(fr,{
 "Traces":"Traces", "Trace exploration":"Exploration des traces", "Exploration mode":"Mode d’exploration", "Observed trace":"Trace observée", "Choose a trace":"Choisir une trace", "spans":"opérations", "Refresh":"Actualiser", "Automatic refresh":"Actualisation automatique", "Trace unavailable":"Trace indisponible", "Revision":"Révision", "Trace view":"Vue de la trace", "tree":"arbre", "timeline":"chronologie", "table":"tableau", "Copy sanitized trace":"Copier la trace filtrée", "Native inclusive span timing":"Durée inclusive native des opérations", "Operation":"Opération", "Producer / service":"Producteur / service", "Start / end (exact ns)":"Début / fin (ns exactes)", "Duration (ns) / quality":"Durée (ns) / qualité", "Partial evidence":"Preuves partielles", "SDK dropped attributes / events / links":"Attributs / événements / liens omis par le SDK", "Projection omitted links":"Liens omis par la projection", "Flags":"Indicateurs", "Typed relationships":"Relations typées", "Relationship":"Relation", "From record":"Enregistrement source", "Target trace / span":"Trace / opération cible", "State":"État", "Span events in export order":"Événements dans l’ordre d’exportation", "Correlated logs":"Journaux corrélés", "resolved":"résolu", "missing":"absent", "ambiguous":"ambigu", "trace-only":"trace seulement", "complete":"complet", "right-censored":"censuré à droite", "unavailable":"indisponible", "unknown":"inconnu",
 "Observed operations and relationships. A trace is not a business journey; inclusive durations are not page dwell or unique total time.":"Opérations et relations observées. Une trace n’est pas un parcours métier ; les durées inclusives ne mesurent ni le temps passé sur la page ni une durée totale sans chevauchement.",
 "Bounded results: narrow collection or inspect a recent trace. Some evidence is omitted.":"Résultats limités : limitez la collecte ou inspectez une trace récente. Certaines preuves sont omises.",
 "No retained traces. Export ordinary spans with valid trace/span IDs; no stage or session attributes are required.":"Aucune trace conservée. Exportez des opérations ordinaires avec des identifiants valides ; aucun attribut d’étape ou de session n’est requis.",
 "Completeness unknown: sampling, missing exports, retention and bounds can omit work. A sampled flag does not establish a sampling rate.":"Exhaustivité inconnue : l’échantillonnage, les exports absents, la conservation et les limites peuvent omettre du travail. Un indicateur d’échantillonnage ne donne pas son taux.",
 "Trace detail truncated at 2,000 records.":"Détail de la trace limité à 2 000 enregistrements."
});

Object.assign(fr,{
 "Give business meaning":"Définir le sens métier", "Business interpretation":"Interprétation métier", "Capture-time definitions":"Définitions à la capture", "Non-persisting preview":"Aperçu sans modification", "capture-time definitions":"Définitions à la capture", "explicit non-persisting preview":"Aperçu explicite sans modification", "Configured goal at span end":"Objectif configuré à la fin de l’opération", "Configured event goal":"Objectif événementiel configuré", "Use this exact mapping as the operational goal":"Utiliser cette correspondance exacte comme objectif opérationnel", "Service (exact; empty means any)":"Service (exact ; vide pour tous)", "Carrier":"Support", "Any carrier":"Tous les supports", "Timing basis":"Base temporelle", "Native span; otherwise unavailable":"Opération native ; sinon indisponible", "Native inclusive operation":"Durée inclusive de l’opération native", "Explicit paired boundaries":"Limites appariées explicites", "Use emitted suggestion":"Utiliser la suggestion émise", "Emitted hint":"Indication émise", "Trace exploration remains available":"L’exploration des traces reste disponible", "A browser intent or generic HTTP success cannot prove persistence. Choose a trusted named commitment or a documented successful operation for that goal.":"Une intention du navigateur ou un succès HTTP générique ne prouve pas l’enregistrement. Pour cet objectif, choisissez un fait nommé fiable après validation ou une opération réussie dont le contrat le garantit."
});

Object.assign(fr,{"Trace references":"Références de traces","No reference":"Aucune référence","Reference name":"Nom de la référence","Save trace reference":"Enregistrer la référence de trace","reference capture-time definitions":"Définitions de capture de la référence","reference pinned definition preview":"Aperçu de la définition enregistrée par la référence"});
Object.assign(fr,{
 'Add app':'Ajouter une application','Choose an app':'Choisir une application','App settings':'Paramètres de l’application','App name':'Nom de l’application','Rename app':'Renommer l’application','Connection settings':'Paramètres de connexion','Register app':'Enregistrer l’application','Reveal credential':'Afficher le jeton','Copy Authorization header':'Copier l’en-tête Authorization','Rotate credential':'Renouveler le jeton','Confirm rotation':'Confirmer le renouvellement','Remove app':'Supprimer l’application','Close':'Fermer','Access unavailable':'Accès indisponible','App unavailable. Choose an app.':'Application indisponible. Choisissez une application.',
 'Configure your existing standard OTLP/HTTP exporter. No Sentrame SDK, app-ID attribute or instrumentation changes are required.':'Configurez votre exportateur OTLP/HTTP standard existant. Aucun SDK Sentrame, attribut d’application ou changement d’instrumentation n’est requis.',
 'Use the reachable base URL for your sender. Signal URLs end in /v1/logs, /v1/traces or /v1/metrics. A base endpoint exporter adds these paths; a signal endpoint uses the full path.':'Utilisez l’URL de base accessible depuis votre émetteur. Les URL des signaux se terminent par /v1/logs, /v1/traces ou /v1/metrics. Un exportateur avec un point de terminaison de base ajoute ces chemins ; un point de terminaison de signal utilise le chemin complet.',
 'Server credentials belong only in server secret configuration. Browser credentials are untrusted; use the supported app-origin forwarder. Direct cross-origin browser export is not qualified.':'Les jetons serveur restent dans la configuration secrète du serveur. Les jetons navigateur ne sont pas fiables ; utilisez le relais pris en charge sur l’origine de l’application. L’exportation directe entre origines n’est pas qualifiée.',
 'Credential rotated. Update the sender configuration.':'Jeton renouvelé. Mettez à jour la configuration de l’émetteur.',
 'Rotation immediately rejects the old credential. Update every sender using this role afterward.':'Le renouvellement rejette immédiatement l’ancien jeton. Mettez ensuite à jour chaque émetteur utilisant ce rôle.',
 'This permanently removes this app and its telemetry, mappings, references and viewer access. Type the exact app name to confirm.':'Cette action supprime définitivement l’application, sa télémétrie, ses correspondances, ses références et ses accès invités. Saisissez son nom exact pour confirmer.',
 'Register your app, then connect its existing OpenTelemetry exporter using the generated settings.':'Enregistrez votre application, puis connectez son exportateur OpenTelemetry existant avec les paramètres générés.',
 'Use a unique name of 1–80 characters. Up to 32 apps per installation.':'Utilisez un nom unique de 1 à 80 caractères. Maximum de 32 applications par installation.',
 'Connect your exporter in App settings, then map meaningful events and supply safe correlation for journeys. Traces remain available.':'Connectez votre exportateur dans les paramètres de l’application, puis définissez des correspondances et une corrélation sûre pour les parcours. Les traces restent disponibles.',
 'use a unique app name of 1–80 characters; maximum 32 apps':'Utilisez un nom unique de 1 à 80 caractères ; maximum de 32 applications.',
 'app registry unavailable; restore a coordinated private backup and restart':'Registre indisponible ; restaurez une sauvegarde privée coordonnée et redémarrez.',
 'Removal pending. Retry removal or restart to resume.':'Suppression en attente. Réessayez ou redémarrez pour la reprendre.','Retry removal':'Réessayer la suppression'
});

Object.assign(fr,{'Copy base endpoint':'Copier le point de terminaison de base','Copy signal endpoint':'Copier le point de terminaison du signal'});

Object.assign(fr,{'App registration unavailable':'Cette inscription de l’application n’est plus disponible.'});

Object.assign(fr,{'secure randomness unavailable':'La génération sécurisée est indisponible.','confirm the exact app name':'Confirmez le nom exact de l’application.'});

Object.assign(fr,{
 'View received records':'Voir les enregistrements reçus',
 'Telemetry has arrived. These records have no derived journey in this selection. Journeys need mapped events and supported correlation; standalone logs remain available as received records.':'La télémétrie a été reçue. Aucun parcours n’est dérivé de ces enregistrements dans cette sélection. Les parcours nécessitent des événements définis et une corrélation prise en charge ; les journaux isolés restent consultables parmi les enregistrements reçus.',
 'No retained traces':'Aucune trace conservée',
 'Logs without trace IDs do not appear in the trace list. This does not mean collection failed. Inspect received records in the journey workspace.':'Les journaux sans identifiant de trace n’apparaissent pas dans la liste des traces. Cela ne signifie pas que la collecte a échoué. Consultez les enregistrements reçus dans l’espace des parcours.'
});

Object.assign(fr,{
 'Enter a stage name before activating.':'Saisissez un nom d’étape avant l’activation.',
 'Choose an operational goal before activating the first mapping.':'Choisissez un objectif opérationnel avant d’activer la première correspondance.',
 'Choose Preview matches before activating. Preview again after changing any field.':'Choisissez Prévisualiser les correspondances avant l’activation. Recommencez après toute modification d’un champ.',
 'Resolve the preview conflict before activating.':'Résolvez le conflit de prévisualisation avant l’activation.',
 'Preview is valid. You can activate this definition.':'La prévisualisation est valide. Vous pouvez activer cette définition.',
 'This app has no operational goal yet. The first definition needs a goal event that ends future journeys. Select the checkbox below only if this is that event.':'Cette application n’a pas encore d’objectif opérationnel. La première définition nécessite un événement objectif qui termine les futurs parcours. Cochez la case ci-dessous uniquement s’il s’agit de cet événement.',
 'End a journey when this event occurs (operational goal)':'Terminer le parcours lorsque cet événement se produit (objectif opérationnel)',
 'Missing journey correlation does not block activation. These records remain unassigned; a mapping alone cannot connect them into journeys.':'L’absence de corrélation de parcours ne bloque pas l’activation. Ces enregistrements restent non attribués ; une correspondance seule ne peut pas les relier en parcours.'
});

Object.assign(fr,{
 'Received records for this mapping':'Enregistrements reçus pour cette correspondance',
 'Examples use retained app-wide history for the selected event, service, carrier and producer, independent of the graph time range.':'Les exemples utilisent l’historique conservé de l’application pour l’événement, le service, le support et le producteur sélectionnés, indépendamment de la période du graphique.',
 'Loading received records…':'Chargement des enregistrements reçus…',
 'before the optional condition':'avant la condition facultative',
 'No retained records match this event context. Choose a received context above or check the exact service, carrier and producer.':'Aucun enregistrement conservé ne correspond à ce contexte. Choisissez un contexte reçu ci-dessus ou vérifiez le service exact, le support et le producteur.',
 'This condition matches no records. The examples and available values below show this event context before applying the condition.':'Cette condition ne correspond à aucun enregistrement. Les exemples et les valeurs ci-dessous montrent ce contexte avant l’application de la condition.',
 'Pick a condition from received data (optional)':'Choisir une condition dans les données reçues (facultatif)',
 'Choose a field and value below, or leave both condition fields empty to include every record in this event context.':'Choisissez un champ et une valeur ci-dessous, ou laissez les deux champs de condition vides pour inclure tous les enregistrements de ce contexte.',
 'No selectable attributes were retained for these records. Leave the condition empty.':'Aucun attribut sélectionnable n’a été conservé pour ces enregistrements. Laissez la condition vide.',
 'Clear condition':'Effacer la condition',
 'Available values are bounded. Some values are omitted.':'Le nombre de valeurs proposées est limité. Certaines valeurs sont omises.',
 'Matching examples':'Exemples correspondants',
 'Examples before the condition':'Exemples avant la condition',
 'Attributes':'Attributs',
 'No retained attributes':'Aucun attribut conservé',
 'No received examples are available.':'Aucun exemple reçu n’est disponible.',
 'Received event context':'Contexte de l’événement reçu',
 'All received contexts':'Tous les contextes reçus',
 'Custom selector':'Sélecteur personnalisé',
 'Use Pick a condition from received data above to fill both fields from an observed value. Only retained, allowed attributes are offered.':'Utilisez Choisir une condition dans les données reçues ci-dessus pour remplir les deux champs avec une valeur observée. Seuls les attributs conservés et autorisés sont proposés.'
});

Object.assign(fr,{'Results are partial. Some retained records were omitted.':'Les résultats sont partiels. Certains enregistrements conservés ont été omis.'});

Object.assign(fr, {
  "A completion goal is optional. Correlated navigation can be explored without one; select this checkbox only to end journeys at this event.": "Un objectif de fin est facultatif. Vous pouvez explorer la navigation corrélée sans objectif ; cochez cette case uniquement pour terminer les parcours à cet événement.",
  "Correlated navigation · no completion goal. Connections show observed order; these paths are excluded from goal conversion.": "Navigation corrélée · aucun objectif de fin. Les liens montrent l’ordre observé ; ces chemins sont exclus du taux de conversion.",
  "Choose 1–100 mappings": "Choisissez de 1 à 100 correspondances",
});

Object.assign(fr, {
  "Completion goal": "Objectif de réussite",
  "Completion goal (optional)": "Objectif de réussite (facultatif)",
  "Set completion goal": "Définir l’objectif de réussite",
  "Current goal": "Objectif actuel",
  "None — navigation only": "Aucun — navigation seulement",
  "Use this event as the completion goal": "Utiliser cet événement comme objectif de réussite",
  "Activate completion goal": "Activer l’objectif de réussite",
  "Event mappings": "Correspondances des événements",
  "A goal is the event that marks a successful journey, such as Meal added to planner. Exit means the journey ended without claiming success. A goal is separate from Kind.": "Un objectif est l’événement qui marque la réussite d’un parcours, comme Repas ajouté au planning. Une sortie termine le parcours sans indiquer de réussite. L’objectif est distinct du type.",
  "This is the current completion goal. To replace it, choose another event in Mappings & goal.": "Il s’agit de l’objectif de réussite actuel. Pour le remplacer, choisissez un autre événement dans Correspondances et objectif.",
  "Activating will make this event the completion goal for new journeys and replace any current goal.": "L’activation définira cet événement comme objectif de réussite des nouveaux parcours et remplacera l’objectif actuel.",
  "Leave unchecked to map this event without changing the completion goal. Journeys can have no goal.": "Laissez cette case décochée pour définir cet événement sans modifier l’objectif de réussite. Les parcours peuvent ne pas avoir d’objectif.",
  "Choose the event that proves the journey succeeded. This is separate from an exit or a dead end. You can explore navigation without setting a goal.": "Choisissez l’événement qui prouve la réussite du parcours. Il est distinct d’une sortie ou d’une impasse. Vous pouvez explorer la navigation sans définir d’objectif.",
  "Choose an event below, review its matching records, then select Preview matches and Activate completion goal.": "Choisissez un événement ci-dessous, examinez les enregistrements correspondants, puis sélectionnez Prévisualiser les correspondances et Activer l’objectif de réussite.",
  "No events received yet. Send an event that represents success, then choose it here.": "Aucun événement reçu pour le moment. Envoyez un événement qui représente une réussite, puis choisissez-le ici.",
  "A completion goal cannot also be an exit. Choose Step for a successful milestone.": "Un objectif de réussite ne peut pas aussi être une sortie. Choisissez Étape pour un jalon de réussite."
});

Object.assign(fr, {
  "Checking this configuration…": "Vérification de la configuration…",
  "Configuration checked. Ready to activate.": "Configuration vérifiée. Prête à être activée.",
  "Could not check this configuration. Try again.": "Impossible de vérifier cette configuration. Réessayez.",
  "Retry check": "Réessayer la vérification",
  "Choose an event below, review its matching records, then select Activate completion goal. Checks update automatically.": "Choisissez un événement ci-dessous, examinez les enregistrements correspondants, puis sélectionnez Activer l’objectif de réussite. Les vérifications se mettent à jour automatiquement.",
  "Leave both fields empty if you do not need this filter. Matching records and validation update automatically before activation.": "Laissez les deux champs vides si ce filtre est inutile. Les enregistrements correspondants et la validation se mettent à jour automatiquement avant l’activation.",
  "matching records are missing a session ID or, for browser events, a tab ID. You can save this mapping, but these records cannot be connected into journeys.": "enregistrements correspondants n’ont pas d’identifiant de session ou, pour les événements du navigateur, d’identifiant d’onglet. Vous pouvez enregistrer cette correspondance, mais ces enregistrements ne peuvent pas être reliés en parcours.",
  "All matching records have the IDs needed to connect them into journeys.": "Tous les enregistrements correspondants ont les identifiants nécessaires pour les relier en parcours."
});

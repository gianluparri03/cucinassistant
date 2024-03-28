from cucinassistant.exceptions import CAError, CACritical
from cucinassistant.database.tests import SubTest
import cucinassistant.database as db


class TestStorage(SubTest):
    def S00_get_storage(self):
        # Tests for get_storage
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.get_storage, self.fake_user)
        self.t.assertEqual(len(db.get_storage(self.giovanna)), 0)
        db.append_storage(self.giovanna, [['one', '2024-03-02', 0], ['two', '2024-03-03', 0], ['three', '2024-03-01', 0]])
        db.append_storage(self.francesco, [['one', '2024-03-02', 0]])
        self.t.assertEqual([a.name for a in db.get_storage(self.giovanna)], ['three', 'one', 'two'])

    def S01_get_storage_article(self):
        # Tests for get_storage_article
        self.t.assertEqual(db.get_storage_article(self.giovanna, 1).name, 'one')
        self.t.assertEqual(db.get_storage_article(self.giovanna, '2').name, 'two')
        self.t.assertRaisesRegex(CAError, 'Articolo non trovato', db.get_storage_article, self.francesco, 1)
        self.t.assertRaisesRegex(CAError, 'Articolo non valido', db.get_storage_article, self.francesco, 'a')

    def S02_append_storage(self):
        # Tests for append_storage
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.append_storage, self.fake_user, [])
        self.t.assertRaisesRegex(CAError, 'Formato non valido', db.append_storage, self.giovanna, [['', '', '', '']])
        self.t.assertRaisesRegex(CAError, 'Nome non valido', db.append_storage, self.giovanna, [['', '', '']])
        self.t.assertRaisesRegex(CAError, 'Quantit&agrave; non valida', db.append_storage, self.giovanna, [['four', '', '2.4']])
        self.t.assertRaisesRegex(CAError, 'Scadenza non valida', db.append_storage, self.giovanna, [['four', '2/4/2023', '']])
        db.append_storage(self.giovanna, [['one', '', '']])
        self.t.assertEqual(db.get_storage_article(self.giovanna, 5), db.Article(5, 'one', None, None))
        db.append_storage(self.giovanna, [['one', '2024-04-02', '']])
        self.t.assertRaisesRegex(CAError, 'Articolo gi&agrave; presente', db.append_storage, self.giovanna, [['one', '2024-03-02', '']])
        self.t.assertEqual(len(db.get_storage(self.giovanna)), 5)
        db.append_storage(self.giovanna, [])

    def S03_remove_storage(self):
        # Tests for remove_storage
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.remove_storage, self.fake_user, [])
        self.t.assertRaisesRegex(CAError, 'Articolo non trovato', db.remove_storage, self.giovanna, [100])
        self.t.assertRaisesRegex(CAError, 'Articolo non trovato', db.remove_storage, self.francesco, [1])
        self.t.assertRaisesRegex(CAError, 'Articolo non valido', db.remove_storage, self.giovanna, ['a'])
        db.remove_storage(self.giovanna, [5, 6])
        self.t.assertEqual(len(db.get_storage(self.giovanna)), 3)
        db.remove_storage(self.giovanna, [])

    def S04_edit_storage(self):
        # Tests for remove_storage
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.edit_storage, self.fake_user, 0, [])
        self.t.assertRaisesRegex(CAError, 'Articolo non valido', db.edit_storage, self.giovanna, 'a', [])
        self.t.assertRaisesRegex(CAError, 'Articolo non trovato', db.edit_storage, self.giovanna, 0, [])
        self.t.assertRaisesRegex(CAError, 'Formato non valido', db.edit_storage, self.giovanna, 1, ['', '', '', ''])
        self.t.assertRaisesRegex(CAError, 'Nome non valido', db.edit_storage, self.giovanna, 1, ['', '', ''])
        self.t.assertRaisesRegex(CAError, 'Quantit&agrave; non valida', db.edit_storage, self.giovanna, 1, ['four', '', '2.4'])
        self.t.assertRaisesRegex(CAError, 'Scadenza non valida', db.edit_storage, self.giovanna, 1, ['four', '2/4/2023', ''])
        db.edit_storage(self.giovanna, 1, ['one', '2024-03-02', 0])
        db.edit_storage(self.giovanna, 1, ['one', '2024-03-02', 3])
        self.t.assertEqual(db.get_storage_article(self.giovanna, 1).quantity, 3)
        self.t.assertRaisesRegex(CAError, 'Articolo gi&agrave; presente', db.edit_storage, self.giovanna, 1, ['two', '2024-03-03', '18'])
        db.edit_storage(self.giovanna, 1, ['two', '2024-03-06', '18'])

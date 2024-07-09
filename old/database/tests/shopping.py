from cucinassistant.exceptions import CAError, CACritical
from cucinassistant.database.tests import SubTest
import cucinassistant.database as db


class TestShopping(SubTest):
    def S00_get_shopping(self):
        # Tests for get_shopping
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.get_shopping, self.fake_user, sec)
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 0)
        db.append_shopping(self.giovanna, sec, ['a', 'b', 'c'])
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 3)

    def S01_append_shopping(self):
        # Tests for append_shopping
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.append_shopping, self.fake_user, sec, [])
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 3)
        db.append_shopping(self.giovanna, sec, ['a', 'b', 'c'])
        db.append_shopping(self.francesco, sec, ['a', 'b', 'c'])
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 3)
        self.t.assertEqual(len(db.get_shopping(self.francesco, sec)), 3)
        db.append_shopping(self.giovanna, sec, [])
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 3)

    def S02_remove_shopping(self):
        # Tests for remove_shopping
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.remove_shopping, self.fake_user, sec, [])
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 3)
        db.remove_shopping(self.giovanna, sec, [e.eid for e in db.get_shopping(self.giovanna, sec) if e.eid > 1])
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 1)
        self.t.assertRaisesRegex(CAError, 'Elemento non trovato', db.remove_shopping, self.giovanna, sec, [100])
        self.t.assertRaisesRegex(CAError, 'Elemento non trovato', db.remove_shopping, self.francesco, sec, [1])
        self.t.assertRaisesRegex(CAError, 'Elemento non valido', db.remove_shopping, self.giovanna, sec, ['a'])
        db.remove_shopping(self.giovanna, sec, ['1'])
        self.t.assertEqual(len(db.get_shopping(self.giovanna, sec)), 0)

    def S03_edit_shopping(self):
        # Tests for remove_shopping
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.edit_shopping, self.fake_user, sec, 0, '')
        self.t.assertRaisesRegex(CAError, 'Elemento non trovato', db.edit_shopping, self.giovanna, sec, 0, '')
        db.append_shopping(self.giovanna, sec, ['before'])
        eid = max(e.eid for e in db.get_shopping(self.giovanna, sec))
        self.t.assertRaisesRegex(CAError, 'Elemento non trovato', db.edit_shopping, self.francesco, sec, eid, '')
        self.t.assertRaisesRegex(CAError, 'Elemento non valido', db.edit_shopping, self.giovanna, sec, 'a', '')
        db.append_shopping(self.giovanna, sec, ['after'])
        db.edit_shopping(self.giovanna, sec, eid, 'before')
        self.t.assertRaisesRegex(CAError, 'Elemento gi&agrave; in lista', db.edit_shopping, self.giovanna, sec, eid, 'after')
        self.t.assertRaisesRegex(CAError, 'Nuovo nome non valido', db.edit_shopping, self.giovanna, sec, eid, '')
        db.edit_shopping(self.giovanna, sec, str(eid), 'after2')
        self.t.assertEqual(db.get_shopping_entry(self.giovanna, sec, eid).name, 'after2')

    def S04_get_shopping_entry(self):
        # Tests for get_shopping_entry
        self.t.assertEqual(db.get_shopping_entry(self.giovanna, sec, 9).name, 'after')
        self.t.assertEqual(db.get_shopping_entry(self.giovanna, sec, '8').name, 'after2')
        self.t.assertRaisesRegex(CAError, 'Elemento non in lista', db.get_shopping_entry, self.francesco, sec, 9)
        self.t.assertRaisesRegex(CAError, 'Elemento non valido', db.get_shopping_entry, self.francesco, sec, 'a')

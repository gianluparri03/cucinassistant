from cucinassistant.exceptions import CAError, CACritical
from cucinassistant.database.tests import SubTest
import cucinassistant.database as db


class TestLists(SubTest):
    def S00_get_list(self):
        # Tests for get_list
        self.t.assertRaisesRegex(CAError, 'Lista inesistente', db.get_list, self.giovanna, '')
        for sec in ('shopping', 'ideas'):
            self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.get_list, self.fake_user, sec)
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 0)
            db.append_list(self.giovanna, sec, ['a', 'b', 'c'])
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 3)

    def S01_append_list(self):
        # Tests for append_list
        self.t.assertRaisesRegex(CAError, 'Lista inesistente', db.append_list, self.giovanna, '', [])
        for sec in ('shopping', 'ideas'):
            self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.append_list, self.fake_user, sec, [])
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 3)
            db.append_list(self.giovanna, sec, ['a', 'b', 'c'])
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 3)
            db.append_list(self.giovanna, sec, [])
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 3)

    def S02_remove_list(self):
        # Tests for remove_list
        self.t.assertRaisesRegex(CAError, 'Lista inesistente', db.remove_list, self.giovanna, '', [])
        for sec in ('shopping', 'ideas'):
            self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.remove_list, self.fake_user, sec, [])
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 3)
            db.remove_list(self.giovanna, sec, [e.eid for e in db.get_list(self.giovanna, sec) if e.eid > 1])
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 1)
            self.t.assertRaisesRegex(CAError, 'Elemento/i non trovato/i', db.remove_list, self.giovanna, sec, [100])
            self.t.assertRaisesRegex(CAError, 'Elemento/i non trovato/i', db.remove_list, self.francesco, sec, [1])
            self.t.assertRaisesRegex(CAError, 'Elemento/i non valido/i', db.remove_list, self.giovanna, sec, ['a'])
            db.remove_list(self.giovanna, sec, ['1'])
            self.t.assertEqual(len(db.get_list(self.giovanna, sec)), 0)

    def S03_edit_list(self):
        # Tests for remove_list
        self.t.assertRaisesRegex(CAError, 'Lista inesistente', db.edit_list, self.giovanna, '', 0, '')
        for sec in ('shopping', 'ideas'):
            self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.edit_list, self.fake_user, sec, 0, '')
            self.t.assertRaisesRegex(CAError, 'Elemento non trovato', db.edit_list, self.giovanna, sec, 0, '')
            db.append_list(self.giovanna, sec, ['before'])
            eid = max(e.eid for e in db.get_list(self.giovanna, sec))
            self.t.assertRaisesRegex(CAError, 'Elemento non trovato', db.edit_list, self.francesco, sec, eid, '')
            self.t.assertRaisesRegex(CAError, 'Elemento non valido', db.edit_list, self.giovanna, sec, 'a', '')
            db.append_list(self.giovanna, sec, ['after'])
            db.edit_list(self.giovanna, sec, eid, 'before')
            self.t.assertRaisesRegex(CAError, 'Elemento gi&agrave; in lista', db.edit_list, self.giovanna, sec, eid, 'after')
            self.t.assertRaisesRegex(CAError, 'Nuovo nome non valido', db.edit_list, self.giovanna, sec, eid, '')
            db.edit_list(self.giovanna, sec, str(eid), 'after2')
            self.t.assertEqual(db.get_list_entry(self.giovanna, sec, eid).name, 'after2')

    def S04_get_list_entry(self):
        # Tests for get_list_entry
        self.t.assertRaisesRegex(CAError, 'Lista inesistente', db.append_list, self.giovanna, '', [])
        for sec in ('shopping', 'ideas'):
            self.t.assertEqual(db.get_list_entry(self.giovanna, sec, 6).name, 'after')
            self.t.assertEqual(db.get_list_entry(self.giovanna, sec, '5').name, 'after2')
            self.t.assertRaisesRegex(CAError, 'Articolo non in lista', db.get_list_entry, self.francesco, sec, 6)
            self.t.assertRaisesRegex(CAError, 'Elemento non valido', db.get_list_entry, self.francesco, sec, 'a')

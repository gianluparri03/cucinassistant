from cucinassistant.exceptions import CAError, CACritical
import cucinassistant.database as db


class TestLists:
    def __init__(self, tester, francesco, giovanna, fake_user):
        self.t = tester
        self.francesco = francesco
        self.giovanna = giovanna
        self.fake_user = fake_user
        self.fake_menu = 0

        # Executes all the subtests IN ORDER
        for name in sorted(dir(self)):
            if name[0] == 'S':
                with self.t.subTest(sec='Lists', sub=int(name[1:3])):
                    getattr(self, name)()


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
            eid = db.append_list(self.giovanna, sec, ['before'])
            self.t.assertRaisesRegex(CAError, 'Elemento non trovato', db.edit_list, self.francesco, sec, eid, '')
            self.t.assertRaisesRegex(CAError, 'Elemento non valido', db.edit_list, self.giovanna, sec, 'a', '')
            db.append_list(self.giovanna, sec, ['after'])
            db.edit_list(self.giovanna, sec, eid, 'before')
            self.t.assertRaisesRegex(CAError, 'Elemento gi&agrave; in lista', db.edit_list, self.giovanna, sec, eid, 'after')
            self.t.assertRaisesRegex(CAError, 'Nuovo nome non valido', db.edit_list, self.giovanna, sec, eid, '')
            db.edit_list(self.giovanna, sec, str(eid), 'after2')
            self.t.assertEqual(db.get_list_entry(self.giovanna, sec, eid), 'after2')

    def S04_get_list_entry(self):
        # Tests for get_list_entry
        self.t.assertRaisesRegex(CAError, 'Lista inesistente', db.append_list, self.giovanna, '', [])
        for sec in ('shopping', 'ideas'):
            self.t.assertEqual(db.get_list_entry(self.giovanna, sec, 6), 'after')
            self.t.assertEqual(db.get_list_entry(self.giovanna, sec, '5'), 'after2')
            self.t.assertRaisesRegex(CAError, 'Articolo non in lista', db.get_list_entry, self.francesco, sec, 6)
            self.t.assertRaisesRegex(CAError, 'Elemento non valido', db.get_list_entry, self.francesco, sec, 'a')

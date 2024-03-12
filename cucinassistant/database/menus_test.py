from cucinassistant.exceptions import CAError, CACritical
import cucinassistant.database as db


class TestMenus:
    def __init__(self, tester, francesco, giovanna, fake_user):
        self.t = tester
        self.francesco = francesco
        self.giovanna = giovanna
        self.fake_user = fake_user
        self.fake_menu = 0

        # Executes all the subtests IN ORDER
        for name in sorted(dir(self)):
            if name[0] == 'S':
                with self.t.subTest(cat='Menus', sub=int(name[1:3])):
                    getattr(self, name)()


    def S00_create_menu(self):
        # Tests for create_menu
        db.create_menu(self.francesco)
        db.create_menu(self.francesco)
        db.create_menu(self.francesco)
        self.t.assertEqual(db.get_menu(self.francesco, 1).next, 2)
        self.t.assertEqual(db.get_menu(self.francesco, 2).prev, 1)
        self.t.assertEqual(db.get_menu(self.francesco, 2).next, 3)
        self.t.assertEqual(db.get_menu(self.francesco, 3).prev, 2)
        self.t.assertEqual(db.get_menu(self.francesco, 3).menu.count(';'), 13)

    def S01_get_menu(self):
        # Tests for get_menu
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.get_menu, self.fake_user)
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.get_menu, self.giovanna, 7)
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.get_menu, self.giovanna, 0)
        self.t.assertEqual(db.get_menu(self.francesco).mid, 3)
        self.t.assertEqual(db.get_menu(self.giovanna).mid, 0)

    def S02_update_menu(self):
        # Tests for update_menu
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.update_menu, self.fake_user, '')
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.update_menu, self.giovanna, 7, '')
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.update_menu, self.giovanna, 0, '')
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non valido', db.update_menu, self.francesco, 1, 'pippo')
        db.update_menu(self.francesco, 1, '0;1;2;3;4;5;6;7;8;9;10;11;12;13')
        self.t.assertEqual(db.get_menu(self.francesco, 1).menu, '0;1;2;3;4;5;6;7;8;9;10;11;12;13')

    def S03_delete_menu(self):
        # Tests for delete_menu
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.update_menu, self.fake_user, '')
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.update_menu, self.giovanna, 7, '')
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.update_menu, self.giovanna, 0, '')
        db.delete_menu(self.francesco, 2)
        self.t.assertEqual(db.get_menu(self.francesco).mid, 3)
        self.t.assertEqual(db.get_menu(self.francesco).prev, 1)
        self.t.assertEqual(db.get_menu(self.francesco, 1).next, 3)
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.get_menu, self.giovanna, 2)

    def S04_duplicate_menu(self):
        # Tests for duplicate_menu
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.duplicate_menu, self.fake_user, '')
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.duplicate_menu, self.giovanna, 7)
        self.t.assertRaisesRegex(CAError, 'Men&ugrave; non trovato', db.duplicate_menu, self.giovanna, 0)
        db.create_menu(self.giovanna)
        db.update_menu(self.giovanna, 1, 'menu;;;;;;;;;;;;;')
        self.t.assertEqual(db.duplicate_menu(self.giovanna, 1), 2)
        self.t.assertEqual(db.get_menu(self.giovanna).mid, 2)
        self.t.assertEqual(db.get_menu(self.giovanna, 2).menu, db.get_menu(self.giovanna, 1).menu)

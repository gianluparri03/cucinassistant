    def S03_delete_menu(self):
        # Tests for delete_menu
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.update_menu, self.fake_user, '')
        self.t.assertRaisesRegex(CAError, 'Menù non trovato', db.update_menu, self.giovanna, 'a', '')
        self.t.assertRaisesRegex(CAError, 'Menù non trovato', db.update_menu, self.giovanna, 7, '')
        self.t.assertRaisesRegex(CAError, 'Menù non trovato', db.update_menu, self.giovanna, 0, '')
        db.delete_menu(self.francesco, 2)
        self.t.assertEqual(db.get_menu(self.francesco).mid, 3)
        self.t.assertEqual(db.get_menu(self.francesco).prev, 1)
        self.t.assertEqual(db.get_menu(self.francesco, 1).next, 3)
        self.t.assertRaisesRegex(CAError, 'Menù non trovato', db.get_menu, self.giovanna, 2)

    def S04_duplicate_menu(self):
        # Tests for duplicate_menu
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.duplicate_menu, self.fake_user, '')
        self.t.assertRaisesRegex(CAError, 'Menù non trovato', db.duplicate_menu, self.giovanna, 'a')
        self.t.assertRaisesRegex(CAError, 'Menù non trovato', db.duplicate_menu, self.giovanna, 7)
        self.t.assertRaisesRegex(CAError, 'Menù non trovato', db.duplicate_menu, self.giovanna, 0)
        db.create_menu(self.giovanna)
        db.update_menu(self.giovanna, 1, 'menu;;;;;;;;;;;;;')
        self.t.assertEqual(db.duplicate_menu(self.giovanna, 1), 2)
        self.t.assertEqual(db.get_menu(self.giovanna).mid, 2)
        self.t.assertEqual(db.get_menu(self.giovanna, 2).menu, db.get_menu(self.giovanna, 1).menu)
